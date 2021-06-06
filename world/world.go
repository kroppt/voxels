package world

import (
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type World struct {
	ubo       *gfx.BufferObject
	cam       *Camera
	chunks    map[ChunkPos]*Chunk
	lastChunk ChunkPos
}

func withinRanges(rng ChunkRange, pos ChunkPos) bool {
	if pos.X < rng.Min.X || pos.X > rng.Max.X {
		return false
	}
	if pos.Z < rng.Min.Z || pos.Z > rng.Max.Z {
		return false
	}
	return true
}

// GetChunkBounds returns the chunk position ranges that are in view around
// currChunk.
func GetChunkBounds(renderDist int32, currChunk ChunkPos) ChunkRange {
	minx := currChunk.X - renderDist
	maxx := currChunk.X + renderDist
	mink := currChunk.Z - renderDist
	maxk := currChunk.Z + renderDist
	return ChunkRange{
		Min: ChunkPos{minx, mink},
		Max: ChunkPos{maxx, maxk},
	}
}

const chunkSize = 10
const chunkHeight = 1
const chunkRenderDist = 2

func NewWorld() (*World, error) {
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

	currChunk := cam.AsVoxelPos().AsChunkPos(chunkSize)
	rng := GetChunkBounds(chunkRenderDist, currChunk)

	chunks := make(map[ChunkPos]*Chunk)
	rng.ForEach(func(pos ChunkPos) {
		ch := NewChunk(chunkSize, chunkHeight, pos)
		chunks[pos] = ch
	})

	world.chunks = chunks
	world.lastChunk = currChunk
	return world, nil
}

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
	key := v.Pos.AsChunkPos(chunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	if chunk, ok := w.chunks[key]; ok {
		chunk.SetVoxel(v)
	}
}

func (w *World) Destroy() {
	w.ubo.Destroy()
}

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

		currChunk := w.cam.AsVoxelPos().AsChunkPos(chunkSize)
		if currChunk != w.lastChunk {
			// the camera position has moved chunks
			// load new chunks
			rng := GetChunkBounds(chunkRenderDist, currChunk)
			rng.ForEach(func(pos ChunkPos) {
				if _, ok := w.chunks[pos]; !ok {
					// chunk i,j is not in map and should be added
					ch := NewChunk(chunkSize, chunkHeight, pos)
					w.chunks[pos] = ch
				}
			})
			// delete old chunks
			lastRng := GetChunkBounds(chunkRenderDist, w.lastChunk)
			lastRng.ForEach(func(pos ChunkPos) {
				inOld := withinRanges(lastRng, pos)
				inNew := withinRanges(rng, pos)
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
