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
	chunks    map[glm.Vec2]*Chunk
	lastChunk glm.Vec2
}

// GetChunkIndex returns the chunk coordinate that the given position
// is in, given the chunkSize.
func GetChunkIndex(chunkSize int, pos glm.Vec3) glm.Vec2 {
	x := int(pos.X())
	z := int(pos.Z())
	if pos.X() < 0 {
		x++
	}
	if pos.Z() < 0 {
		z++
	}
	x /= chunkSize
	z /= chunkSize
	if pos.X() < 0 {
		x--
	}
	if pos.Z() < 0 {
		z--
	}
	return glm.Vec2{float32(x), float32(z)}
}

type Range struct {
	Min int
	Max int
}

func withinRanges(xrng, yrng Range, x, z int) bool {
	if x < xrng.Min || x > xrng.Max {
		return false
	}
	if z < yrng.Min || z > yrng.Max {
		return false
	}
	return true
}

func applyWithinRanges(xrng, zrng Range, fn func(i, j int)) {
	for x := xrng.Min; x <= xrng.Max; x++ {
		for z := zrng.Min; z <= zrng.Max; z++ {
			fn(x, z)
		}
	}
}

// GetChunkBounds returns the minimum and maximum chunks indices that should be
// in view around a camera at the given chunk.
func GetChunkBounds(worldSize int, currChunk glm.Vec2) (x Range, z Range) {
	halfWorld := worldSize / 2
	minx := int(currChunk.X()) - halfWorld
	maxx := int(currChunk.X()) + halfWorld
	minj := int(currChunk.Y()) - halfWorld
	maxj := int(currChunk.Y()) + halfWorld
	return Range{minx, maxx}, Range{minj, maxj}
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

	currChunk := GetChunkIndex(chunkSize, cam.GetPosition())
	xrng, zrng := GetChunkBounds(chunkRenderDist*2+1, currChunk)

	chunks := make(map[glm.Vec2]*Chunk)
	applyWithinRanges(xrng, zrng, func(i, j int) {
		ch := NewChunk(chunkSize, chunkHeight, glm.Vec2{float32(i), float32(j)})
		chunks[glm.Vec2{float32(i), float32(j)}] = ch
	})

	world.chunks = chunks
	world.lastChunk = currChunk
	return world, nil
}

func (w *World) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
	var hits []*Voxel
	for _, chunk := range w.chunks {
		chunkHits, _ := chunk.root.Find(func(node *Octree) bool {
			half := node.GetAABC().Size / float32(2.0)
			aabc := AABC{
				Pos:  (&node.GetAABC().Pos).Add(&glm.Vec3{half, half, half}),
				Size: node.GetAABC().Size,
			}
			_, hit := Intersect(aabc, w.cam.GetPosition(), w.cam.GetLookForward())
			return hit
		})
		hits = append(hits, chunkHits...)
	}

	closest, dist := GetClosest(w.cam.GetPosition(), hits)
	return closest, dist, len(hits) != 0
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

		currChunk := GetChunkIndex(chunkSize, w.cam.GetPosition())
		if currChunk != w.lastChunk {
			// the camera position has moved chunks
			// load new chunks
			xrng, zrng := GetChunkBounds(chunkRenderDist*2+1, currChunk)
			applyWithinRanges(xrng, zrng, func(i, j int) {
				key := glm.Vec2{float32(i), float32(j)}
				if _, ok := w.chunks[key]; !ok {
					// chunk i,j is not in map and should be added
					ch := NewChunk(chunkSize, chunkHeight, glm.Vec2{float32(i), float32(j)})
					w.chunks[key] = ch
				}
			})
			// delete old chunks
			lastXRange, lastZRange := GetChunkBounds(chunkRenderDist*2+1, w.lastChunk)
			applyWithinRanges(lastXRange, lastZRange, func(i, j int) {
				oldKey := glm.Vec2{float32(i), float32(j)}
				inOld := withinRanges(lastXRange, lastZRange, i, j)
				inNew := withinRanges(xrng, zrng, i, j)
				if inOld && !inNew {
					w.chunks[oldKey].Destroy()
					delete(w.chunks, oldKey)
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
