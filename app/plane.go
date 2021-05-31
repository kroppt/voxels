package app

import (
	"errors"
	"fmt"

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
	col := [3]float32{vox.Color.R, vox.Color.G, vox.Color.B}
	colors := [8][3]float32{
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
}

func NewPlaneRenderer() *PlaneRenderer {
	return &PlaneRenderer{}
}

func (pr *PlaneRenderer) Destroy() {
	cleanupCubes(pr.cubes)
}

func (pr *PlaneRenderer) Init(plane *world.Plane) error {
	x, y, z := plane.Size()
	cubes, err := makeCubes(plane, x, y, z)
	if err != nil {
		return fmt.Errorf("failed to initialize cubes: %v", err)
	}
	pr.plane = plane
	pr.cubes = cubes
	return nil
}

var errWrongPlanePassed = errors.New("wrong Plane was passed into Render")

func (pr *PlaneRenderer) Render(plane *world.Plane) error {
	if pr.plane != plane {
		return errWrongPlanePassed
	}
	for _, xcubes := range pr.cubes {
		for _, ycubes := range xcubes {
			for _, cube := range ycubes {
				cube.Render(*plane.GetCamera())
			}
		}
	}
	return nil
}
