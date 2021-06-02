package world

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/shapes"
)

type Range struct {
	Min int
	Max int
}

type Plane struct {
	x      Range
	y      Range
	z      Range
	voxels [][][]*Voxel
	ubo    *gfx.BufferObject
	cam    *Camera
}

func makeVoxel(x, y, z Range, i, j, k int) (*Voxel, error) {
	pos := Position{X: i, Y: j, Z: k}
	r := float32(i-x.Min) / float32(x.Max-x.Min)
	g := float32(j-y.Min) / float32(y.Max-y.Min)
	b := float32(k-z.Min) / float32(z.Max-z.Min)
	a := float32(1.0)
	col := [4]float32{r, g, b, a}
	colors := [8][4]float32{
		col, col, col, col, col, col, col, col,
	}
	obj, err := shapes.NewColoredCube(colors)
	if err != nil {
		return nil, fmt.Errorf("couldn't create colored cube at %v: %w", pos, err)
	}
	obj.Translate(float32(i), float32(j), float32(k))
	return &Voxel{
		Object: obj,
	}, nil
}

func makeYVoxels(x, y, z Range, i, j int) ([]*Voxel, error) {
	yvox := []*Voxel{}
	for k := z.Min; k <= z.Max; k++ {
		zvox, err := makeVoxel(x, y, z, i, j, k)
		if err != nil {
			return yvox, err
		}
		yvox = append(yvox, zvox)
	}
	return yvox, nil
}

func makeXVoxels(x, y, z Range, i int) ([][]*Voxel, error) {
	xvox := [][]*Voxel{}
	for j := y.Min; j <= y.Max; j++ {
		yvox, err := makeYVoxels(x, y, z, i, j)
		if err != nil {
			return xvox, err
		}
		xvox = append(xvox, yvox)
	}
	return xvox, nil
}

func makeVoxels(x, y, z Range) ([][][]*Voxel, error) {
	voxels := [][][]*Voxel{}
	for i := x.Min; i <= x.Max; i++ {
		xvox, err := makeXVoxels(x, y, z, i)
		if err != nil {
			return voxels, err
		}
		voxels = append(voxels, xvox)
	}
	return voxels, nil
}

func cleanupCubes(objs [][][]*Voxel) {
	for _, objs1 := range objs {
		for _, objs2 := range objs1 {
			for _, obj := range objs2 {
				obj.Destroy()
			}
		}
	}
}

func NewPlane(x, y, z Range) (*Plane, error) {
	voxels, err := makeVoxels(x, y, z)
	if err != nil {
		cleanupCubes(voxels)
		return nil, err
	}

	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)

	cam := NewCamera()

	plane := &Plane{
		x:      x,
		y:      y,
		z:      z,
		voxels: voxels,
		ubo:    ubo,
		cam:    cam,
	}
	plane.cam.SetPosition(&glm.Vec3{0, 0, 25})
	plane.cam.LookAt(&glm.Vec3{0, 0, 0})
	err = plane.UpdateView()
	if err != nil {
		return nil, err
	}
	err = plane.UpdateProj()
	if err != nil {
		return nil, err
	}
	return plane, nil
}

func (p *Plane) Destroy() {
	p.ubo.Destroy()
}

var ErrOutOfBounds = errors.New("position out of bounds of plane")

func getRangeOffsets(pos Position, x, y, z Range) (i int, j int, k int) {
	i = pos.X - x.Min
	j = pos.Y - y.Min
	k = pos.Z - z.Min
	return i, j, k
}

func (p *Plane) At(pos Position) (*Voxel, error) {
	switch {
	case pos.X < p.x.Min:
	case pos.X > p.x.Max:
	case pos.Y < p.y.Min:
	case pos.Y > p.y.Max:
	case pos.Z < p.z.Min:
	case pos.Z > p.z.Max:
	default:
		i, j, k := getRangeOffsets(pos, p.x, p.y, p.z)
		return p.voxels[i][j][k], nil
	}
	return nil, ErrOutOfBounds
}

func (p *Plane) Size() (x, y, z Range) {
	return p.x, p.y, p.z
}

func (p *Plane) GetCamera() *Camera {
	return p.cam
}

func (p *Plane) UpdateView() error {
	cam := *p.GetCamera()
	view := cam.GetViewMat()
	err := p.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	return nil
}

func (p *Plane) UpdateProj() error {
	cam := *p.GetCamera()
	proj := cam.GetProjMat()
	err := p.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

var ErrNilRenderer = errors.New("cannot render with nil renderer")

func (p *Plane) Render() {
	for _, xcubes := range p.voxels {
		for _, ycubes := range xcubes {
			for _, cube := range ycubes {
				cube.Render()
			}
		}
	}
}
