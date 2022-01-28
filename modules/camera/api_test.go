package camera_test

import (
	"fmt"
	"math"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := camera.New(nil, nil, player.PositionEvent{})
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestModuleNewInitialPos(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		expected := player.PositionEvent{
			X: 1,
			Y: 2,
			Z: 3,
		}
		var actual player.PositionEvent
		playerMod := &player.FnModule{
			FnUpdatePlayerPosition: func(posEvent player.PositionEvent) {
				actual = posEvent
			},
		}
		camera.New(playerMod, nil, player.PositionEvent{
			X: 1,
			Y: 2,
			Z: 3,
		})
		if actual != expected {
			t.Fatalf("expected initial camera pos %v but got %v", expected, actual)
		}
	})
}

func TestMovementEventPlayerMod(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		direction camera.MoveDirection
		x         int32
		y         int32
		z         int32
	}{
		{
			direction: camera.MoveForwards,
			x:         0,
			y:         0,
			z:         -1,
		},
		{
			direction: camera.MoveRight,
			x:         1,
			y:         0,
			z:         0,
		},
		{
			direction: camera.MoveBackwards,
			x:         0,
			y:         0,
			z:         1,
		},
		{
			direction: camera.MoveLeft,
			x:         -1,
			y:         0,
			z:         0,
		},
		{
			direction: 0,
			x:         0,
			y:         0,
			z:         0,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("updates chunk module position when moving %v", tC.direction.String()), func(t *testing.T) {
			t.Parallel()

			var actualEvent player.PositionEvent
			playerMod := &player.FnModule{
				FnUpdatePlayerPosition: func(posEvent player.PositionEvent) {
					actualEvent = posEvent
				},
			}
			graphicsMod := &graphics.FnModule{}

			mod := camera.New(playerMod, graphicsMod, player.PositionEvent{})

			movementEvent := camera.MovementEvent{
				Direction: tC.direction,
			}
			mod.HandleMovementEvent(movementEvent)

			expectEvent := player.PositionEvent{
				X: tC.x,
				Y: tC.y,
				Z: tC.z,
			}
			if actualEvent != expectEvent {
				t.Fatalf("expected %v but got %v", expectEvent, actualEvent)
			}
		})
	}

	t.Run("updates chunk module position when moving multiple times", func(t *testing.T) {
		t.Parallel()

		var actualEvent player.PositionEvent
		playerMod := &player.FnModule{
			FnUpdatePlayerPosition: func(posEvent player.PositionEvent) {
				actualEvent = posEvent
			},
		}
		graphicsMod := &graphics.FnModule{}

		mod := camera.New(playerMod, graphicsMod, player.PositionEvent{})

		moveRightEvent := camera.MovementEvent{
			Direction: camera.MoveRight,
		}
		moveBackEvent := camera.MovementEvent{
			Direction: camera.MoveBackwards,
		}
		mod.HandleMovementEvent(moveRightEvent)
		mod.HandleMovementEvent(moveBackEvent)

		expectEvent := player.PositionEvent{
			X: 1,
			Y: 0,
			Z: 1,
		}
		if actualEvent != expectEvent {
			t.Fatalf("expected %v but got %v", expectEvent, actualEvent)
		}
	})
}

func withinError(x, y float64, diff float64) bool {
	if x+diff > y && x-diff < y {
		return true
	}
	return false
}

func withinErrorVec3(a, b mgl.Vec3, diff float64) bool {
	if withinError(a.X(), b.X(), diff) && withinError(a.Y(), b.Y(), diff) &&
		withinError(a.Z(), b.Z(), diff) {
		return true
	}
	return false
}

func withinErrorQuat(q1 mgl.Quat, q2 mgl.Quat, diff float64) bool {
	if withinError(q1.W, q2.W, diff) && withinErrorVec3(q1.V, q2.V, diff) {
		return true
	}
	return false
}

const errMargin = 0.000001

func TestModuleHandleLookEvent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		looks []struct {
			right float64
			down  float64
		}
		rotation mgl.Quat
	}{
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 1.0,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float64{0, -math.Sin(1.0 / 2), 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: -1.0,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float64{0, math.Sin(1.0 / 2), 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 0.0,
					down:  1.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float64{-math.Sin(1.0 / 2), 0, 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 0.0,
					down:  -1.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float64{math.Sin(1.0 / 2), 0, 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 3.0 * math.Pi / 4.0,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(3.0 * math.Pi / 8.0),
				V: [3]float64{0.0, -math.Sin(3.0 * math.Pi / 8.0), 0.0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 0.0,
					down:  -math.Pi / 4.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(math.Pi / 8.0),
				V: [3]float64{math.Sin(math.Pi / 8.0), 0.0, 0.0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: math.Pi,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: 0,
				V: [3]float64{0, -1, 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: math.Pi / 2.0,
					down:  math.Pi / 2.0,
				},
			},
			rotation: mgl.Quat{
				W: 1.0 / 2.0,
				V: [3]float64{-1.0 / 2.0, -1.0 / 2.0, -1.0 / 2.0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: math.Pi / 2.0,
					down:  math.Pi / 2.0,
				},
				{
					right: math.Pi / 2.0,
					down:  -math.Pi / 2.0,
				},
			},
			rotation: mgl.Quat{
				W: 0,
				V: [3]float64{0, -1, 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: 0.0,
					down:  math.Pi / 2.0,
				},
				{
					right: math.Pi / 2.0,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: 1.0 / 2.0,
				V: [3]float64{-1.0 / 2.0, -1.0 / 2.0, -1.0 / 2.0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: -math.Pi / 4.0,
					down:  -math.Pi / 4.0,
				},
				{
					right: math.Pi / 4.0,
					down:  0.0,
				},
			},
			rotation: mgl.Quat{
				W: math.Cos(math.Pi / 8.0),
				V: [3]float64{math.Sin(math.Pi / 8.0), 0, 0},
			},
		},
		{
			looks: []struct {
				right float64
				down  float64
			}{
				{
					right: -math.Pi / 4.0,
					down:  -math.Pi / 4.0,
				},
				{
					right: math.Pi / 4.0,
					down:  0.0,
				},
				{
					right: 0,
					down:  math.Pi / 4.0,
				},
			},
			rotation: mgl.QuatIdent(),
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("updates graphics module direction when looking %+v", tC.looks), func(t *testing.T) {
			t.Parallel()

			var evt graphics.DirectionEvent
			graphicsMod := graphics.FnModule{
				FnUpdatePlayerDirection: func(directionEvent graphics.DirectionEvent) {
					evt = directionEvent
				},
			}

			mod := camera.New(nil, graphicsMod, player.PositionEvent{})

			for _, look := range tC.looks {
				lookEvent := camera.LookEvent{
					Right: look.right,
					Down:  look.down,
				}
				mod.HandleLookEvent(lookEvent)
			}

			expected := graphics.DirectionEvent{
				Rotation: tC.rotation,
			}
			if !withinErrorQuat(expected.Rotation, evt.Rotation, errMargin) {
				t.Fatalf("expected %v but got %v", expected, evt)
			}
		})
	}
}
