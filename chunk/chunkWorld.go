package chunk

import (
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type ChunkWorld struct {
	ubo    *gfx.BufferObject
	cam    *Camera
	chunks map[glm.Vec2]*Chunk
}

// GetChunkAt returns the chunk coordinate that the given position is in, given
// the chunkSize.
func GetChunkAt(chunkSize int, pos glm.Vec3) glm.Vec2 {
	x := int(pos.X()) / chunkSize
	z := int(pos.Z()) / chunkSize
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

// chunkWorldSize must be odd.
const chunkWorldSize = 3
const chunkSize = 6
const chunkHeight = 1

func NewChunkWorld() (*ChunkWorld, error) {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	cam := NewCamera()
	world := &ChunkWorld{
		ubo: ubo,
		cam: cam,
	}
	cam.SetPosition(&glm.Vec3{3, 2, 3})
	cam.LookAt(&glm.Vec3{0, 0, 0})

	currChunk := GetChunkAt(chunkSize, cam.GetPosition())
	xrng, zrng := GetChunkBounds(chunkWorldSize, currChunk)

	chunks := make(map[glm.Vec2]*Chunk)
	for i := xrng.Min; i <= xrng.Max; i++ {
		for k := zrng.Min; k <= zrng.Max; k++ {
			ch := NewChunk(chunkSize, chunkHeight, glm.Vec2{float32(i), float32(k)})
			chunks[glm.Vec2{float32(i), float32(k)}] = ch
		}
	}
	world.chunks = chunks

	return world, nil
}

func (w *ChunkWorld) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
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

func (w *ChunkWorld) Destroy() {
	w.ubo.Destroy()
}

func (w *ChunkWorld) GetCamera() *Camera {
	return w.cam
}

func (w *ChunkWorld) updateView() error {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	cam.Clean()
	return nil
}

func (w *ChunkWorld) updateProj() error {
	cam := *w.GetCamera()
	proj := cam.GetProjMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	cam.Clean()
	return nil
}

func (w *ChunkWorld) Render() error {
	if w.cam.IsDirty() {
		err := w.updateView()
		if err != nil {
			return err
		}
		err = w.updateProj()
		if err != nil {
			return err
		}
		currChunk := GetChunkAt(chunkSize, w.cam.GetPosition())
		xrng, zrng := GetChunkBounds(chunkWorldSize, currChunk)
		for i := xrng.Min; i <= xrng.Max; i++ {
			for j := zrng.Min; j <= zrng.Max; j++ {
				key := glm.Vec2{float32(i), float32(j)}
				if _, ok := w.chunks[key]; !ok {
					ch := NewChunk(chunkSize, chunkHeight, glm.Vec2{float32(i), float32(j)})
					w.chunks[key] = ch
				}
			}
		}
	}
	for _, chunk := range w.chunks {
		chunk.Render()
	}
	return nil
}
