package player_test

import (
	"fmt"
	"math"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := player.New(nil, nil)
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func testMovementEventChunkMod(t *testing.T) {
	testCases := []struct {
		direction player.MoveDirection
		x         int32
		y         int32
		z         int32
	}{
		{
			direction: player.MoveForwards,
			x:         0,
			y:         0,
			z:         -1,
		},
		{
			direction: player.MoveRight,
			x:         1,
			y:         0,
			z:         0,
		},
		{
			direction: player.MoveBackwards,
			x:         0,
			y:         0,
			z:         1,
		},
		{
			direction: player.MoveLeft,
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

			var actualEvent chunk.PositionEvent
			chunkMod := &chunk.FnModule{
				FnUpdatePlayerPosition: func(posEvent chunk.PositionEvent) {
					actualEvent = posEvent
				},
			}
			graphicsMod := &graphics.FnModule{}

			mod := player.New(chunkMod, graphicsMod)

			movementEvent := player.MovementEvent{
				Direction: tC.direction,
			}
			mod.HandleMovementEvent(movementEvent)

			expectEvent := chunk.PositionEvent{
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

		var actualEvent chunk.PositionEvent
		chunkMod := &chunk.FnModule{
			FnUpdatePlayerPosition: func(posEvent chunk.PositionEvent) {
				actualEvent = posEvent
			},
		}
		graphicsMod := &graphics.FnModule{}

		mod := player.New(chunkMod, graphicsMod)

		moveRightEvent := player.MovementEvent{
			Direction: player.MoveRight,
		}
		moveBackEvent := player.MovementEvent{
			Direction: player.MoveBackwards,
		}
		mod.HandleMovementEvent(moveRightEvent)
		mod.HandleMovementEvent(moveBackEvent)

		expectEvent := chunk.PositionEvent{
			X: 1,
			Y: 0,
			Z: 1,
		}
		if actualEvent != expectEvent {
			t.Fatalf("expected %v but got %v", expectEvent, actualEvent)
		}
	})
}

func testMovementEventGraphicsMod(t *testing.T) {
	testCases := []struct {
		name        string
		moveEvents  []player.MovementEvent
		expectEvent graphics.PositionEvent
	}{
		{
			name: "forward 1",
			moveEvents: []player.MovementEvent{
				{Direction: player.MoveForwards, Pressed: true},
				{Direction: player.MoveForwards, Pressed: false},
			},
			expectEvent: graphics.PositionEvent{
				X: 0,
				Y: 0,
				Z: -1,
			},
		},
		{
			name: "right 1",
			moveEvents: []player.MovementEvent{
				{Direction: player.MoveRight, Pressed: true},
				{Direction: player.MoveRight, Pressed: false},
			},
			expectEvent: graphics.PositionEvent{
				X: 1,
				Y: 0,
				Z: 0,
			},
		},
		{
			name: "back 1",
			moveEvents: []player.MovementEvent{
				{Direction: player.MoveBackwards, Pressed: true},
				{Direction: player.MoveBackwards, Pressed: false},
			},
			expectEvent: graphics.PositionEvent{
				X: 0,
				Y: 0,
				Z: 1,
			},
		},
		{
			name: "left 1",
			moveEvents: []player.MovementEvent{
				{Direction: player.MoveLeft, Pressed: true},
				{Direction: player.MoveLeft, Pressed: false},
			},
			expectEvent: graphics.PositionEvent{
				X: -1,
				Y: 0,
				Z: 0,
			},
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run("updates graphics module position when moving "+tC.name, func(t *testing.T) {
			t.Parallel()

			chunkMod := &chunk.FnModule{}
			var actualEvent *graphics.PositionEvent
			graphicsMod := &graphics.FnModule{
				FnUpdatePlayerPosition: func(posEvent graphics.PositionEvent) {
					actualEvent = &posEvent
				},
			}

			mod := player.New(chunkMod, graphicsMod)

			for _, me := range tC.moveEvents {
				mod.HandleMovementEvent(me)
			}

			if actualEvent == nil {
				t.Fatal("expected event but got none")
			}

			if *actualEvent != tC.expectEvent {
				t.Fatalf("expected %v but got %v", tC.expectEvent, *actualEvent)
			}
		})
	}
}

func TestModuleHandleMovementEvent(t *testing.T) {
	t.Parallel()

	testMovementEventChunkMod(t)

	testMovementEventGraphicsMod(t)
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

			mod := player.New(nil, graphicsMod)

			for _, look := range tC.looks {
				lookEvent := player.LookEvent{
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
