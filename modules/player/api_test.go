package player_test

import (
	"fmt"
	"testing"

	"github.com/EngoEngine/math"
	"github.com/engoengine/glm"
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

func TestModuleHandleLookEvent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		right    float32
		down     float32
		rotation glm.Quat
	}{
		{
			right: 1.0,
			down:  0.0,
			rotation: glm.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float32{0, -math.Sin(1.0 / 2), 0},
			},
		},
		{
			right: -1.0,
			down:  0.0,
			rotation: glm.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float32{0, math.Sin(1.0 / 2), 0},
			},
		},
		{
			right: 0.0,
			down:  1.0,
			rotation: glm.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float32{-math.Sin(1.0 / 2), 0, 0},
			},
		},
		{
			right: 0.0,
			down:  -1.0,
			rotation: glm.Quat{
				W: math.Cos(1.0 / 2),
				V: [3]float32{math.Sin(1.0 / 2), 0, 0},
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("updates graphics module direction when looking %v right and %v down", tC.right, tC.down), func(t *testing.T) {
			t.Parallel()

			var evt graphics.DirectionEvent
			graphicsMod := &fnGraphicsMod{
				fnUpdateDirection: func(directionEvent graphics.DirectionEvent) {
					evt = directionEvent
				},
			}

			mod := player.New(nil, graphicsMod)

			lookEvent := player.LookEvent{
				Right: tC.right,
				Down:  tC.down,
			}
			mod.HandleLookEvent(lookEvent)

			expected := graphics.DirectionEvent{
				Rotation: tC.rotation,
			}
			if evt != expected {
				t.Fatalf("expected %v but got %v", expected, evt)
			}
		})
	}
}