package world_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/engoengine/glm"
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
		_ = cam.GetPosition()
	})

	t.Run("initial camera position is 0, 0, 0", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		pos := cam.GetPosition()
		expect := glm.Vec3{0.0, 0.0, 0.0}

		if pos != expect {
			t.Fatalf("expected %v but got %v", expect, pos)
		}
	})

}

func TestTableCameraTranslate(t *testing.T) {

	testCases := []struct {
		start  glm.Vec3
		diff   glm.Vec3
		expect glm.Vec3
	}{
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{0.0, 0.0, 0.0},
			expect: glm.Vec3{0.0, 0.0, 0.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{1.0, 0.0, 0.0},
			expect: glm.Vec3{1.0, 0.0, 0.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{1.0, 1.0, 0.0},
			expect: glm.Vec3{1.0, 1.0, 0.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{1.0, 1.0, 1.0},
			expect: glm.Vec3{1.0, 1.0, 1.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{-1.0, 0.0, 0.0},
			expect: glm.Vec3{-1.0, 0.0, 0.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{-1.0, -1.0, 0.0},
			expect: glm.Vec3{-1.0, -1.0, 0.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{-1.0, -1.0, -1.0},
			expect: glm.Vec3{-1.0, -1.0, -1.0},
		},
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{25.0, 12.3, 0.1},
			expect: glm.Vec3{25.0, 12.3, 0.1},
		},
		{
			start:  glm.Vec3{10.0, 15.0, 20.0},
			diff:   glm.Vec3{-10.0, -15.0, -20.0},
			expect: glm.Vec3{0.0, 0.0, 0.0},
		},
		{
			start:  glm.Vec3{5.0, 5.0, 5.0},
			diff:   glm.Vec3{-10.0, -10.0, -10.0},
			expect: glm.Vec3{-5.0, -5.0, -5.0},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("translate %v by %v", tC.start, tC.diff), func(t *testing.T) {
			t.Parallel()
			cam := world.NewCamera()
			cam.Translate(&tC.start)

			cam.Translate(&tC.diff)
			pos := cam.GetPosition()

			if pos != tC.expect {
				t.Fatalf("expected %v but got %v", tC.expect, pos)
			}
		})
	}

}

func TestCameraGetRotationQuat(t *testing.T) {

	t.Run("get camera rotation", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		_ = cam.GetRotationQuat()
	})

	t.Run("initial camera rotation is identity quat", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		rot := cam.GetRotationQuat()
		expect := glm.QuatIdent()

		if rot != expect {
			t.Fatalf("expected %v but got %v", expect, rot)
		}
	})

}

func TestTableCameraRotate(t *testing.T) {

	testCases := []struct {
		start  glm.Vec3
		diff   glm.Vec3
		expect glm.Quat
	}{
		{
			start:  glm.Vec3{0.0, 0.0, 0.0},
			diff:   glm.Vec3{0.0, 0.0, 0.0},
			expect: glm.QuatIdent(),
		},
		{
			start: glm.Vec3{0.0, 0.0, 0.0},
			diff:  glm.Vec3{90.0, 0.0, 0.0},
			expect: glm.Quat{
				W: float32(math.Cos(float64(glm.DegToRad(90.0 / 2)))),
				V: glm.Vec3{
					float32(math.Sin(float64(glm.DegToRad(90.0 / 2)))),
					0.0,
					0.0,
				},
			},
		},
		{
			start: glm.Vec3{90.0, 0.0, 0.0},
			diff:  glm.Vec3{0.0, 90.0, 0.0},
			expect: (&glm.Quat{
				W: float32(math.Cos(float64(glm.DegToRad(90.0 / 2)))),
				V: glm.Vec3{
					float32(math.Sin(float64(glm.DegToRad(90.0 / 2)))),
					0.0,
					0.0,
				},
			}).Mul(&glm.Quat{
				W: float32(math.Cos(float64(glm.DegToRad(90.0 / 2)))),
				V: glm.Vec3{
					0.0,
					float32(math.Sin(float64(glm.DegToRad(90.0 / 2)))),
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
			cam.Rotate(&glm.Vec3{1.0, 0.0, 0.0}, tC.start.X())
			cam.Rotate(&glm.Vec3{0.0, 1.0, 0.0}, tC.start.Y())
			cam.Rotate(&glm.Vec3{0.0, 0.0, 1.0}, tC.start.Z())

			cam.Rotate(&glm.Vec3{1.0, 0.0, 0.0}, tC.diff.X())
			cam.Rotate(&glm.Vec3{0.0, 1.0, 0.0}, tC.diff.Y())
			cam.Rotate(&glm.Vec3{0.0, 0.0, 1.0}, tC.diff.Z())
			quat := cam.GetRotationQuat()
			margin := float32(0.00001)
			if !withinError(quat.W, tC.expect.W, margin) || !withinErrorVec3(quat.V, tC.expect.V, margin) {
				t.Fatalf("expected %v but got %v", tC.expect, quat)
			}
		})
	}

}

func withinError(x, y float32, diff float32) bool {
	if x+diff > y && x-diff < y {
		return true
	}
	return false
}

func withinErrorVec3(a, b glm.Vec3, diff float32) bool {
	if withinError(a.X(), b.X(), diff) && withinError(a.Y(), b.Y(), diff) &&
		withinError(a.Z(), b.Z(), diff) {
		return true
	}
	return false
}

func TestCameraLookVector(t *testing.T) {
	t.Run("roll doesn't change look forward", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		cam.Rotate(&glm.Vec3{0, 0, 1}, 60)
		look := cam.GetLookForward()
		expect := glm.Vec3{0.0, 0.0, -1.0}
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})
	t.Run("back is 180 deg yaw from forward", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		expect := cam.GetLookBack()
		cam.Rotate(&glm.Vec3{0, 1, 0}, 180)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})
	t.Run("right is +90 deg yaw from forward", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		expect := cam.GetLookRight()
		cam.Rotate(&glm.Vec3{0, 1, 0}, 90)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})
	t.Run("left is -90 deg yaw from forward", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		expect := cam.GetLookLeft()
		cam.Rotate(&glm.Vec3{0, 1, 0}, -90)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})

	t.Run("up is -90 deg pitch from forward", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		expect := cam.GetLookUp()
		cam.Rotate(&glm.Vec3{1, 0, 0}, -90)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})

	t.Run("fancy rotation", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		// Use the right hand rule:
		// Point thumb in direction of axis, fingers curl towards
		// negative degree rotation
		expect := cam.GetLookUp()
		cam.Rotate(&glm.Vec3{1, 0, 0}, 90)  // look down
		cam.Rotate(&glm.Vec3{0, 1, 0}, 90)  // roll 90 toward right
		cam.Rotate(&glm.Vec3{0, 0, 1}, -90) // look up (facing right)
		cam.Rotate(&glm.Vec3{1, 0, 0}, 180) // roll upsidedown
		cam.Rotate(&glm.Vec3{0, 0, 1}, 270) // lean back 270 degrees (now facing up)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})

	t.Run("rotate about negative axis", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		// Right hand rule works for negative axes too
		expect := cam.GetLookDown()
		cam.Rotate(&glm.Vec3{-1, 0, 0}, -90)
		look := cam.GetLookForward()
		if !withinErrorVec3(look, expect, 0.0001) {
			t.Fatalf("expected %v but got %v", expect, look)
		}
	})

}

func TestCameraVoxelCoords(t *testing.T) {
	testCases := []struct {
		desc   string
		cPos   glm.Vec3
		expect world.VoxelPos
	}{
		{
			desc:   "pos pos works",
			cPos:   glm.Vec3{0.5, 0.5, 0.5},
			expect: world.VoxelPos{0, 0, 0},
		},
		{
			desc:   "pos neg works",
			cPos:   glm.Vec3{0.5, 0.5, -0.5},
			expect: world.VoxelPos{0, 0, -1},
		},
		{
			desc:   "neg pos works",
			cPos:   glm.Vec3{-0.5, 0.5, 0.5},
			expect: world.VoxelPos{-1, 0, 0},
		},
		{
			desc:   "neg neg works",
			cPos:   glm.Vec3{-0.5, 0.5, -0.5},
			expect: world.VoxelPos{-1, 0, -1},
		},
		{
			desc:   "far neg neg works",
			cPos:   glm.Vec3{-1.5, -2.5, -1.5},
			expect: world.VoxelPos{-2, -3, -2},
		},
		{
			desc:   "far pos pos works",
			cPos:   glm.Vec3{1.5, 2.5, 1.5},
			expect: world.VoxelPos{1, 2, 1},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cam := world.NewCamera()
			cam.SetPosition(&tC.cPos)
			result := cam.AsVoxelPos()
			if result != tC.expect {
				t.Fatalf("expected cam voxel coord %v but got %v", tC.expect, result)
			}
		})
	}
}
