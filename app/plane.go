package app

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/shapes"
	"github.com/kroppt/voxels/voxgl"
	"github.com/kroppt/voxels/world"
)

func makeCube(plane *world.Plane, i, j, k int) (*voxgl.Object, error) {
	x, y, z := plane.Size()
	pos := world.Position{X: i, Y: j, Z: k}
	vox, err := plane.At(pos)
	if err != nil {
		return nil, fmt.Errorf("couldn't get voxel at %v: %w", pos, err)
	}
	vox.Color.R = float32(i-x.Min) / float32(x.Max-x.Min)
	vox.Color.G = float32(j-y.Min) / float32(y.Max-y.Min)
	vox.Color.B = float32(k-z.Min) / float32(z.Max-z.Min)
	vox.Color.A = 1.0
	col := [4]float32{vox.Color.R, vox.Color.G, vox.Color.B, vox.Color.A}
	colors := [8][4]float32{
		col, col, col, col, col, col, col, col,
	}
	obj, err := shapes.NewColoredCube(colors)
	if err != nil {
		return nil, fmt.Errorf("couldn't create colored cube at %v: %w", pos, err)
	}
	obj.Translate(float32(i), float32(j), float32(k))
	return obj, err
}

func makeYCubes(plane *world.Plane, i, j int, z world.Range) ([]*voxgl.Object, error) {
	ycube := []*voxgl.Object{}
	for k := z.Min; k <= z.Max; k++ {
		zcube, err := makeCube(plane, i, j, k)
		if err != nil {
			return nil, err
		}
		ycube = append(ycube, zcube)
	}
	return ycube, nil
}

func makeXCubes(plane *world.Plane, i int, y, z world.Range) ([][]*voxgl.Object, error) {
	xcube := [][]*voxgl.Object{}
	for j := y.Min; j <= y.Max; j++ {
		ycube, err := makeYCubes(plane, i, j, z)
		if err != nil {
			return nil, err
		}
		xcube = append(xcube, ycube)
	}
	return xcube, nil
}

func cleanupCubes(objs [][][]*voxgl.Object) {
	for _, objs1 := range objs {
		for _, objs2 := range objs1 {
			for _, obj := range objs2 {
				obj.Destroy()
			}
		}
	}
}

func makeCubes(plane *world.Plane, x, y, z world.Range) ([][][]*voxgl.Object, error) {
	cubes := [][][]*voxgl.Object{}
	for i := x.Min; i <= x.Max; i++ {
		xcube, err := makeXCubes(plane, i, y, z)
		if err != nil {
			cleanupCubes(cubes)
			return nil, err
		}
		cubes = append(cubes, xcube)
	}
	return cubes, nil
}

type PlaneRenderer struct {
	plane *world.Plane
	cubes [][][]*voxgl.Object
	ubo   *gfx.BufferObject
}

func NewPlaneRenderer() *PlaneRenderer {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	return &PlaneRenderer{
		plane: nil,
		cubes: nil,
		ubo:   ubo,
	}
}

func (pr *PlaneRenderer) Destroy() {
	cleanupCubes(pr.cubes)
	pr.ubo.Destroy()
}

func (pr *PlaneRenderer) Init(plane *world.Plane) error {
	x, y, z := plane.Size()
	cubes, err := makeCubes(plane, x, y, z)
	if err != nil {
		return fmt.Errorf("failed to initialize cubes: %v", err)
	}
	pr.plane = plane
	pr.cubes = cubes
	err = pr.UpdateView()
	if err != nil {
		return err
	}
	err = pr.UpdateProj()
	if err != nil {
		return err
	}
	return nil
}

func (pr *PlaneRenderer) UpdateView() error {
	cam := *pr.plane.GetCamera()
	view := cam.GetViewMat()
	err := pr.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	return nil
}

func (pr *PlaneRenderer) UpdateProj() error {
	cam := *pr.plane.GetCamera()
	proj := cam.GetProjMat()
	err := pr.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

var errWrongPlanePassed = errors.New("wrong Plane was passed into Render")

func (pr *PlaneRenderer) Render(plane *world.Plane) error {
	if pr.plane != plane {
		return errWrongPlanePassed
	}

	cam := *plane.GetCamera()
	for _, xcubes := range pr.cubes {
		for _, ycubes := range xcubes {
			for _, cube := range ycubes {
				cube.Render(cam)
			}
		}
	}
	return nil
}
