package world

import (
	"errors"
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Range struct {
	Min int
	Max int
}

type PlaneRenderer interface {
	Init(*Plane) error
	Render(*Plane) error
}

type Plane struct {
	renderer PlaneRenderer
	x        Range
	y        Range
	z        Range
	voxels   [][][]Voxel
	cam      *Camera
}

func makeVoxel(i, j, k int) Voxel {
	return Voxel{
		Position: Position{i, j, k},
		Color:    Color{0.0, 0.0, 0.0, 0.0},
	}
}

func makeYVoxels(i, j int, z Range) []Voxel {
	yvox := []Voxel{}
	for k := z.Min; k <= z.Max; k++ {
		zvox := makeVoxel(i, j, k)
		yvox = append(yvox, zvox)
	}
	return yvox
}

func makeXVoxels(i int, y, z Range) [][]Voxel {
	xvox := [][]Voxel{}
	for j := y.Min; j <= y.Max; j++ {
		yvox := makeYVoxels(i, j, z)
		xvox = append(xvox, yvox)
	}
	return xvox
}

func makeVoxels(x, y, z Range) [][][]Voxel {
	voxels := [][][]Voxel{}
	for i := x.Min; i <= x.Max; i++ {
		xvox := makeXVoxels(i, y, z)
		voxels = append(voxels, xvox)
	}
	return voxels
}

func NewPlane(renderer PlaneRenderer, x, y, z Range) (*Plane, error) {
	voxels := makeVoxels(x, y, z)
	plane := &Plane{
		renderer: renderer,
		x:        x,
		y:        y,
		z:        z,
		voxels:   voxels,
		cam:      NewCamera(),
	}
	plane.cam.Translate(mgl.Vec3{0.0, 0.0, -25.0})
	plane.cam.Rotate(mgl.Vec3{1.0, 0.0, 0.0}, 45.0)
	plane.cam.Rotate(mgl.Vec3{0.0, 1.0, 0.0}, 45.0)

	if renderer != nil {
		err := renderer.Init(plane)
		if err != nil {
			return nil, fmt.Errorf("could not initialize plane renderer: %w", err)
		}
	}
	return plane, nil
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
		return &p.voxels[i][j][k], nil
	}
	return nil, ErrOutOfBounds
}

func (p *Plane) Size() (x, y, z Range) {
	return p.x, p.y, p.z
}

func (p *Plane) GetCamera() *Camera {
	return p.cam
}

var ErrNilRenderer = errors.New("cannot render with nil renderer")

func (p *Plane) Render() error {
	if p.renderer == nil {
		return ErrNilRenderer
	}
	return p.renderer.Render(p)
}
