package world

import (
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type ChunkWorld struct {
	ubo    *gfx.BufferObject
	cam    *Camera
	chunks []*Chunk
}

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
	world.cam.SetPosition(&glm.Vec3{0, 0, 25})
	world.cam.LookAt(&glm.Vec3{0, 0, 0})

	// TODO implement me

	return world, nil
}

func (w *ChunkWorld) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
	// TODO loop over all chunks

	// candidates, ok := w.root.Find(func(node *Octree) bool {
	// 	half := node.GetAABC().Size / float32(2.0)
	// 	aabc := AABC{
	// 		Pos:  (&node.GetAABC().Pos).Add(&glm.Vec3{half, half, half}),
	// 		Size: node.GetAABC().Size,
	// 	}
	// 	_, hit := Intersect(aabc, w.cam.GetPosition(), w.cam.GetLookForward())
	// 	return hit
	// })
	// closest, dist := GetClosest(w.cam.GetPosition(), candidates)
	// return closest, dist, ok
	return nil, 0, false
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
	}
	// TODO do rendering
	return nil
}
