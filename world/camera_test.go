package world_test

import (
	"fmt"
	"math"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/kroppt/voxels/world"
)

func TestNewCamera(t *testing.T) {

	t.Run("new camera not nil", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		if cam == nil {
			t.Fatal("expected valid camera but got nil")
		}
	})

}

func TestCameraGetPosition(t *testing.T) {

	t.Run("get camera position", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		if cam == nil {
			t.Fatal("expected valid camera but got nil")
		}
		_ = cam.GetPosition()
	})

	t.Run("initial camera position is 0, 0, 0", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		if cam == nil {
			t.Fatal("expected valid camera but got nil")
		}
		pos := cam.GetPosition()
		expect := mgl.Vec3{0.0, 0.0, 0.0}

		if pos != expect {
			t.Fatalf("expected %v but got %v", expect, pos)
		}
	})

}

func TestTableCameraTranslate(t *testing.T) {

	testCases := []struct {
		start  mgl.Vec3
		diff   mgl.Vec3
		expect mgl.Vec3
	}{
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{0.0, 0.0, 0.0},
			expect: mgl.Vec3{0.0, 0.0, 0.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{1.0, 0.0, 0.0},
			expect: mgl.Vec3{1.0, 0.0, 0.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{1.0, 1.0, 0.0},
			expect: mgl.Vec3{1.0, 1.0, 0.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{1.0, 1.0, 1.0},
			expect: mgl.Vec3{1.0, 1.0, 1.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{-1.0, 0.0, 0.0},
			expect: mgl.Vec3{-1.0, 0.0, 0.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{-1.0, -1.0, 0.0},
			expect: mgl.Vec3{-1.0, -1.0, 0.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{-1.0, -1.0, -1.0},
			expect: mgl.Vec3{-1.0, -1.0, -1.0},
		},
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{25.0, 12.3, 0.1},
			expect: mgl.Vec3{25.0, 12.3, 0.1},
		},
		{
			start:  mgl.Vec3{10.0, 15.0, 20.0},
			diff:   mgl.Vec3{-10.0, -15.0, -20.0},
			expect: mgl.Vec3{0.0, 0.0, 0.0},
		},
		{
			start:  mgl.Vec3{5.0, 5.0, 5.0},
			diff:   mgl.Vec3{-10.0, -10.0, -10.0},
			expect: mgl.Vec3{-5.0, -5.0, -5.0},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("translate %v by %v", tC.start, tC.diff), func(t *testing.T) {
			t.Parallel()
			cam := world.NewCamera()
			if cam == nil {
				t.Fatal("expected valid camera but got nil")
			}
			cam.Translate(tC.start)

			cam.Translate(tC.diff)
			pos := cam.GetPosition()

			if pos != tC.expect {
				t.Fatalf("expected %v but got %v", tC.expect, pos)
			}
		})
	}

}

func TestCameraGetRotation(t *testing.T) {

	t.Run("get camera rotation", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		if cam == nil {
			t.Fatal("expected valid camera but got nil")
		}
		_ = cam.GetRotation()
	})

	t.Run("initial camera rotation is identity quat", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		if cam == nil {
			t.Fatal("expected valid camera but got nil")
		}
		rot := cam.GetRotation()
		expect := mgl.QuatIdent()

		if rot != expect {
			t.Fatalf("expected %v but got %v", expect, rot)
		}
	})

}

func TestTableCameraRotate(t *testing.T) {

	testCases := []struct {
		start  mgl.Vec3
		diff   mgl.Vec3
		expect mgl.Quat
	}{
		{
			start:  mgl.Vec3{0.0, 0.0, 0.0},
			diff:   mgl.Vec3{0.0, 0.0, 0.0},
			expect: mgl.QuatIdent(),
		},
		{
			start: mgl.Vec3{0.0, 0.0, 0.0},
			diff:  mgl.Vec3{90.0, 0.0, 0.0},
			expect: mgl.Quat{
				W: float32(math.Cos(float64(mgl.DegToRad(90.0 / 2)))),
				V: mgl.Vec3{
					float32(math.Sin(float64(mgl.DegToRad(90.0 / 2)))),
					0.0,
					0.0,
				},
			},
		},
		{
			start: mgl.Vec3{90.0, 0.0, 0.0},
			diff:  mgl.Vec3{0.0, 90.0, 0.0},
			expect: mgl.Quat{
				W: float32(math.Cos(float64(mgl.DegToRad(90.0 / 2)))),
				V: mgl.Vec3{
					float32(math.Sin(float64(mgl.DegToRad(90.0 / 2)))),
					0.0,
					0.0,
				},
			}.Mul(mgl.Quat{
				W: float32(math.Cos(float64(mgl.DegToRad(90.0 / 2)))),
				V: mgl.Vec3{
					0.0,
					float32(math.Sin(float64(mgl.DegToRad(90.0 / 2)))),
					0.0,
				},
			}),
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("rotate %v by %v", tC.start, tC.diff), func(t *testing.T) {
			t.Parallel()
			cam := world.NewCamera()
			if cam == nil {
				t.Fatal("expected valid camera but got nil")
			}
			cam.Rotate(mgl.Vec3{1.0, 0.0, 0.0}, tC.start.X())
			cam.Rotate(mgl.Vec3{0.0, 1.0, 0.0}, tC.start.Y())
			cam.Rotate(mgl.Vec3{0.0, 0.0, 1.0}, tC.start.Z())

			cam.Rotate(mgl.Vec3{1.0, 0.0, 0.0}, tC.diff.X())
			cam.Rotate(mgl.Vec3{0.0, 1.0, 0.0}, tC.diff.Y())
			cam.Rotate(mgl.Vec3{0.0, 0.0, 1.0}, tC.diff.Z())
			quat := cam.GetRotation()

			if quat != tC.expect {
				t.Fatalf("expected %v but got %v", tC.expect, quat)
			}
		})
	}

}
