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

type fnChunkMod struct {
	fnUpdatePosition func(chunk.PositionEvent)
}

func (fn fnChunkMod) UpdatePosition(posEvent chunk.PositionEvent) {
	fn.fnUpdatePosition(posEvent)
}

type fnGraphicsMod struct {
	fnUpdateDirection func(graphics.DirectionEvent)
}

func (fn fnGraphicsMod) UpdateDirection(directionEvent graphics.DirectionEvent) {
	fn.fnUpdateDirection(directionEvent)
}

func TestModuleHandleMovementEvent(t *testing.T) {
	t.Parallel()

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

			var evt chunk.PositionEvent
			chunkMod := &fnChunkMod{
				fnUpdatePosition: func(posEvent chunk.PositionEvent) {
					evt = posEvent
				},
			}

			mod := player.New(chunkMod, nil)

			movementEvent := player.MovementEvent{
				Direction: tC.direction,
			}
			mod.HandleMovementEvent(movementEvent)

			expected := chunk.PositionEvent{
				X: tC.x,
				Y: tC.y,
				Z: tC.z,
			}
			if evt != expected {
				t.Fatalf("expected %v but got %v", expected, evt)
			}
		})
	}

	t.Run("updates chunk module position when moving multiple times", func(t *testing.T) {
		var evt chunk.PositionEvent
		chunkMod := &fnChunkMod{
			fnUpdatePosition: func(posEvent chunk.PositionEvent) {
				evt = posEvent
			},
		}

		mod := player.New(chunkMod, nil)

		moveRightEvent := player.MovementEvent{
			Direction: player.MoveRight,
		}
		moveBackEvent := player.MovementEvent{
			Direction: player.MoveBackwards,
		}
		mod.HandleMovementEvent(moveRightEvent)
		mod.HandleMovementEvent(moveBackEvent)

		expected := chunk.PositionEvent{
			X: 1,
			Y: 0,
			Z: 1,
		}
		if evt != expected {
			t.Fatalf("expected %v but got %v", expected, evt)
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
			graphicsMod := &fnGraphicsMod{
				fnUpdateDirection: func(directionEvent graphics.DirectionEvent) {
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
