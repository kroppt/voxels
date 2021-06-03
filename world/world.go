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
		Object:      obj,
		coordinates: glm.Vec3{float32(i), float32(j), float32(k)},
	}, nil
}

// func makeYVoxels(x, y, z Range, i, j int) ([]*Voxel, error) {
// 	yvox := []*Voxel{}
// 	for k := z.Min; k <= z.Max; k++ {
// 		zvox, err := makeVoxel(x, y, z, i, j, k)
// 		if err != nil {
// 			return yvox, err
// 		}
// 		yvox = append(yvox, zvox)
// 	}
// 	return yvox, nil
// }

// func makeXVoxels(x, y, z Range, i int) ([][]*Voxel, error) {
// 	xvox := [][]*Voxel{}
// 	for j := y.Min; j <= y.Max; j++ {
// 		yvox, err := makeYVoxels(x, y, z, i, j)
// 		if err != nil {
// 			return xvox, err
// 		}
// 		xvox = append(xvox, yvox)
// 	}
// 	return xvox, nil
// }

// func makeVoxels(x, y, z Range) ([][][]*Voxel, error) {
// 	voxels := [][][]*Voxel{}
// 	for i := x.Min; i <= x.Max; i++ {
// 		xvox, err := makeXVoxels(x, y, z, i)
// 		if err != nil {
// 			return voxels, err
// 		}
// 		voxels = append(voxels, xvox)
// 	}
// 	return voxels, nil
// }

// func cleanupCubes(objs [][][]*Voxel) {
// 	for _, objs1 := range objs {
// 		for _, objs2 := range objs1 {
// 			for _, obj := range objs2 {
// 				obj.Destroy()
// 			}
// 		}
// 	}
// }

func makeOctree(x, y, z Range) (tree *Octree, err error) {
	for i := x.Min; i <= x.Max; i++ {
		for j := y.Min; j <= y.Max; j++ {
			for k := z.Min; k <= z.Max; k++ {
				voxel, err := makeVoxel(x, y, z, i, j, k)
				if err != nil {
					return tree, err
				}
				tree.addLeaf(voxel)
			}
		}
	}
	return tree, nil
}

func cleanupOctree(tree *Octree) {
	curr := tree.children
	for curr != nil {
		curr.node.voxel.Destroy()
		curr = curr.next
	}
	if tree.voxel != nil {
		tree.voxel.Destroy()
	}
}

func New(x, y, z Range) (*World, error) {
	octree, err := makeOctree(x, y, z)
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

func (p *World) Destroy() {
	p.ubo.Destroy()
}

var ErrOutOfBounds = errors.New("position out of bounds")

func getRangeOffsets(pos Position, x, y, z Range) (i int, j int, k int) {
	i = pos.X - x.Min
	j = pos.Y - y.Min
	k = pos.Z - z.Min
	return i, j, k
}

func (w *World) At(pos Position) (*Voxel, error) {
	switch {
	case pos.X < w.x.Min:
	case pos.X > w.x.Max:
	case pos.Y < w.y.Min:
	case pos.Y > w.y.Max:
	case pos.Z < w.z.Min:
	case pos.Z > w.z.Max:
	default:
		i, j, k := getRangeOffsets(pos, w.x, w.y, w.z)
		return w.voxels[i][j][k], nil
	}
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
	for _, xcubes := range w.voxels {
		for _, ycubes := range xcubes {
			for _, cube := range ycubes {
				cube.Render()
			}
		}
	}
}
