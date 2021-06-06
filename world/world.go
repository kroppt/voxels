package world

import (
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// World tracks the camera and its renderable chunks.
type World struct {
	ubo       *gfx.BufferObject
	cam       *Camera
	chunks    map[ChunkPos]*Chunk
	lastChunk ChunkPos
}

const chunkSize = 10
const chunkHeight = 1
const chunkRenderDist = 2

// New returns a new world.World.
func New() *World {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	cam := NewCamera()
	world := &World{
		ubo: ubo,
		cam: cam,
	}
	cam.SetPosition(&glm.Vec3{3, 2, 3})
	cam.LookAt(&glm.Vec3{0, 0, 0})

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
	var hits []*Voxel
	for _, chunk := range w.chunks {
		chunkHits, _ := chunk.root.Find(func(node *Octree) bool {
			aabc := *node.GetAABC()
			_, hit := Intersect(aabc, w.cam.GetPosition(), w.cam.GetLookForward())
			return hit
		})
		hits = append(hits, chunkHits...)
	}

	closest, dist := GetClosest(w.cam.GetPosition(), hits)
	return closest, dist, len(hits) != 0
}

// SetVoxel updates a voxel's variables in the world if the chunk
// that it would belong to is currently loaded.
func (w *World) SetVoxel(v *Voxel) {
	key := v.Pos.GetChunkPos(chunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	if chunk, ok := w.chunks[key]; ok {
		chunk.SetVoxel(v)
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

func (w *World) updateView() error {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	return nil
}

func (w *World) updateProj() error {
	cam := *w.GetCamera()
	proj := cam.GetProjMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

// Render renders the chunks of the world in OpenGL.
// TODO isolate chunk loading and unloading logic.
func (w *World) Render() error {
	if w.cam.IsDirty() {
		err := w.updateView()
		if err != nil {
			return err
		}
		err = w.updateProj()
		if err != nil {
			return err
		}
		w.cam.Clean()

		currChunk := w.cam.AsVoxelPos().GetChunkPos(chunkSize)
		if currChunk != w.lastChunk {
			// the camera position has moved chunks
			// load new chunks
			rng := currChunk.GetSurroundings(chunkRenderDist)
			rng.ForEach(func(pos ChunkPos) {
				if _, ok := w.chunks[pos]; !ok {
					// chunk i,j is not in map and should be added
					ch := NewChunk(chunkSize, chunkHeight, pos)
					w.chunks[pos] = ch
				}
			})
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
	for _, chunk := range w.chunks {
		chunk.Render()
	}
	return nil
}
