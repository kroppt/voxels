package world

import (
	"fmt"
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/util"
)

// World tracks the camera and its renderable chunks.
type World struct {
	ubo       *gfx.BufferObject
	cam       *Camera
	chunks    map[ChunkPos]*Chunk
	lastChunk ChunkPos
}

const chunkSize = 16
const chunkHeight = 5
const chunkRenderDist = 10

// New returns a new world.World.
func New() *World {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	cam := NewCameraDefault()

	world := &World{
		ubo: ubo,
		cam: cam,
	}

	cam.SetPosition(&glm.Vec3{0.5, 0.5, 2})
	cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})
	currChunk := cam.AsVoxelPos().GetChunkPos(chunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)

	chunks := make(map[ChunkPos]*Chunk)
	rng.ForEach(func(pos ChunkPos) {
		ch := NewChunk(chunkSize, chunkHeight, pos)
		chunks[pos] = ch
	})

	world.chunks = chunks
	world.lastChunk = currChunk
	return world
}

// FindLookAtVoxel determines which voxel is being looked at. It returns the
// block, distance to the block, and whether the block was found.
func (w *World) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
	var bestVox *Voxel
	var bestDist float32
	for _, chunk := range w.chunks {
		vox, dist, hit := chunk.root.FindClosestIntersect(w.cam)
		if hit && (dist < bestDist || bestVox == nil) {
			bestVox = vox
			bestDist = dist
		}
	}
	return bestVox, bestDist, bestVox != nil
}

// SetVoxel updates a voxel's variables in the world if the chunk
// that it would belong to is currently loaded.
func (w *World) SetVoxel(v *Voxel) {
	key := v.Pos.GetChunkPos(chunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	if chunk, ok := w.chunks[key]; ok {
		chunk.SetVoxel(v, 0x3F)
	}
}

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
}

// GetCamera returns a reference to the camera.
func (w *World) GetCamera() *Camera {
	return w.cam
}

// TODO Create UBO extraction
func (w *World) updateUBO() error {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	proj := cam.GetProjMat()
	err = w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(view)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

// Render renders the chunks of the world in OpenGL.
// TODO isolate chunk loading and unloading logic.
func (w *World) Render() error {
	if w.cam.IsDirty() {
		err := w.updateUBO()
		if err != nil {
			return err
		}
		w.cam.Clean()

		currChunk := w.cam.AsVoxelPos().GetChunkPos(chunkSize)
		if currChunk != w.lastChunk {
			sw := util.Start()
			// the camera position has moved chunks
			// load new chunks
			rng := currChunk.GetSurroundings(chunkRenderDist)
			counter := 0
			rng.ForEach(func(pos ChunkPos) {
				if _, ok := w.chunks[pos]; !ok {
					counter++
					// chunk i,j is not in map and should be added
					ch := NewChunk(chunkSize, chunkHeight, pos)
					w.chunks[pos] = ch
				}
			})
			sw.StopRecordAverage(fmt.Sprintf("Load %v Chunks", counter))
			// delete old chunks
			lastRng := w.lastChunk.GetSurroundings(chunkRenderDist)
			lastRng.ForEach(func(pos ChunkPos) {
				inOld := lastRng.Contains(pos)
				inNew := rng.Contains(pos)
				if inOld && !inNew {
					w.chunks[pos].Destroy()
					delete(w.chunks, pos)
				}
			})
			w.lastChunk = currChunk
		}
	}
	culled := 0
	for _, chunk := range w.chunks {
		if w.cam.IsWithinFrustum(chunk.Pos.AsVec3(), float32(chunk.size), float32(chunk.height), float32(chunk.size)) {
			chunk.Render()
		} else {
			culled++
		}
	}
	// log.Debugf("culled %v / %v chunks", culled, len(w.chunks))
	return nil
}
