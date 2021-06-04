package world

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/voxgl"
)

type Range struct {
	Min int
	Max int
}

type World struct {
	x    Range
	y    Range
	z    Range
	ubo  *gfx.BufferObject
	cam  *Camera
	root *Octree
}

func makeVoxel(x, y, z Range, i, j, k int) (*Voxel, error) {
	pos := Position{X: i, Y: j, Z: k}
	r := float32(i-x.Min) / float32(x.Max-x.Min)
	g := float32(j-y.Min) / float32(y.Max-y.Min)
	b := float32(k-z.Min) / float32(z.Max-z.Min)
	a := float32(1.0)
	point := [7]float32{float32(i), float32(j), float32(k), r, g, b, a}
	obj, err := voxgl.NewColoredObject(point)
	if err != nil {
		return nil, fmt.Errorf("couldn't create colored object at %v: %w", pos, err)
	}

	return &Voxel{
		Object: obj,
		Pos:    glm.Vec3{float32(i), float32(j), float32(k)},
	}, nil
}

func fillOctree(x, y, z Range) (tree *Octree, err error) {
	for i := x.Min; i <= x.Max; i++ {
		for j := y.Min; j <= y.Max; j++ {
			for k := z.Min; k <= z.Max; k++ {
				voxel, err := makeVoxel(x, y, z, i, j, k)
				if err != nil {
					return tree, err
				}
				tree = tree.AddLeaf(voxel)
			}
		}
	}
	return tree, nil
}

func cleanupOctree(tree *Octree) {
	tree.Apply(func(o *Octree) {
		o.voxel.Destroy()
	})
}

func New(x, y, z Range) (*World, error) {
	octree, err := fillOctree(x, y, z)
	if err != nil {
		cleanupOctree(octree)
		return nil, err
	}

	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)

	cam := NewCamera()

	world := &World{
		x:    x,
		y:    y,
		z:    z,
		ubo:  ubo,
		cam:  cam,
		root: octree,
	}
	world.cam.SetPosition(&glm.Vec3{0, 0, 25})
	world.cam.LookAt(&glm.Vec3{0, 0, 0})
	err = world.UpdateView()
	if err != nil {
		cleanupOctree(octree)
		return nil, err
	}
	err = world.UpdateProj()
	if err != nil {
		cleanupOctree(octree)
		return nil, err
	}
	return world, nil
}

func (w *World) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
	candidates, ok := w.root.Find(func(node *Octree) bool {
		half := node.GetAABC().Size / float32(2.0)
		aabc := AABC{
			Pos:  (&node.GetAABC().Pos).Add(&glm.Vec3{half, half, half}),
			Size: node.GetAABC().Size,
		}
		_, hit := Intersect(aabc, w.cam.GetPosition(), w.cam.GetLookForward())
		return hit
	})
	closest, dist := GetClosest(w.cam.GetPosition(), candidates)
	return closest, dist, ok
}

func (w *World) Destroy() {
	w.ubo.Destroy()
	cleanupOctree(w.root)
}

var ErrOutOfBounds = errors.New("position out of bounds")

func (w *World) At(pos Position) (*Voxel, error) {
	// TODO unimplemented, do we want this?
	return nil, ErrOutOfBounds
}

func (w *World) Size() (x, y, z Range) {
	return w.x, w.y, w.z
}

func (w *World) GetCamera() *Camera {
	return w.cam
}

func (w *World) UpdateView() error {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	return nil
}

func (w *World) UpdateProj() error {
	cam := *w.GetCamera()
	proj := cam.GetProjMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

func (w *World) Render() {
	w.root.Apply(func(o *Octree) {
		o.voxel.Render()
	})
}
