package camera_test

import (
	"fmt"
	"math"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/player"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := camera.New(&player.FnModule{}, player.PositionEvent{})
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestModuleNewInitialPos(t *testing.T) {
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
	camera.New(playerMod, player.PositionEvent{
		X: 1,
		Y: 2,
		Z: 3,
	})
	if actual != expected {
		t.Fatalf("expected initial camera pos %v but got %v", expected, actual)
	}
}

func TestOnlyHandlingMovementEventDoesNotMovePlayer(t *testing.T) {
	t.Parallel()
	expected := false
	actual := false
	playerMod := &player.FnModule{}
	cameraMod := camera.New(playerMod, player.PositionEvent{})
	playerMod.FnUpdatePlayerPosition = func(posEvent player.PositionEvent) {
		actual = true
	}

	cameraMod.HandleMovementEvent(camera.MovementEvent{
		Direction: camera.MoveForwards,
		Pressed:   true,
	})
	if actual != expected {
		t.Fatal("expected player's position to not be updated, but it was")
	}
}

func TestCameraTickOnlyUpdatesPlayerPositionIfMoved(t *testing.T) {
	t.Parallel()
	expected := false
	actual := false
	playerMod := &player.FnModule{}
	cameraMod := camera.New(playerMod, player.PositionEvent{})
	playerMod.FnUpdatePlayerPosition = func(posEvent player.PositionEvent) {
		actual = true
	}
	cameraMod.Tick()
	if actual != expected {
		t.Fatal("expected player's position to not be updated, but it was")
	}
}

func TestCameraDoesNotMoveIfTickOccursWhileNoMovementKeyIsPressed(t *testing.T) {
	t.Parallel()
	playerMod := &player.FnModule{}
	cameraMod := camera.New(playerMod, player.PositionEvent{})
	expected := false
	actual := false
	playerMod.FnUpdatePlayerPosition = func(posEvent player.PositionEvent) {
		actual = true
	}
	cameraMod.HandleMovementEvent(camera.MovementEvent{
		Direction: camera.MoveForwards,
		Pressed:   true,
	})
	cameraMod.HandleMovementEvent(camera.MovementEvent{
		Direction: camera.MoveForwards,
		Pressed:   false,
	})
	cameraMod.Tick()
	if actual != expected {
		t.Fatal("expected player's position to not be updated, but it was")
	}
}

func TestCameraMovesIfTickOccursWhileNoMovementKeyIsPressed(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		direction camera.MoveDirection
		x         float64
		y         float64
		z         float64
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
		t.Run(fmt.Sprintf("updates player module position when moving %v", tC.direction.String()), func(t *testing.T) {
			t.Parallel()

			var actualEvent player.PositionEvent
			playerMod := &player.FnModule{
				FnUpdatePlayerPosition: func(posEvent player.PositionEvent) {
					actualEvent = posEvent
				},
			}
			cameraMod := camera.New(playerMod, player.PositionEvent{})

			movementEvent := camera.MovementEvent{
				Direction: tC.direction,
				Pressed:   true,
			}
			cameraMod.HandleMovementEvent(movementEvent)
			cameraMod.Tick()

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

	t.Run("updates player module position when moving multiple times", func(t *testing.T) {
		t.Parallel()

		var actualEvent player.PositionEvent
		playerMod := &player.FnModule{
			FnUpdatePlayerPosition: func(posEvent player.PositionEvent) {
				actualEvent = posEvent
			},
		}

		cameraMod := camera.New(playerMod, player.PositionEvent{})

		moveRightEvent := camera.MovementEvent{
			Direction: camera.MoveRight,
			Pressed:   true,
		}
		moveBackEvent := camera.MovementEvent{
			Direction: camera.MoveBackwards,
			Pressed:   true,
		}
		cameraMod.HandleMovementEvent(moveRightEvent)
		cameraMod.HandleMovementEvent(moveBackEvent)
		cameraMod.Tick()

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

func TestPlayerMovesInDirectionOfCamera(t *testing.T) {
	testCases := []struct {
		desc     string
		lookEvt  camera.LookEvent
		expected mgl.Vec3
	}{
		{
			desc: "move forward after looking right",
			lookEvt: camera.LookEvent{
				Right: math.Pi / 2.0,
				Down:  0,
			},
			expected: mgl.Vec3{1, 0, 0},
		},
		{
			desc: "move forward after looking left",
			lookEvt: camera.LookEvent{
				Right: -math.Pi / 2.0,
				Down:  0,
			},
			expected: mgl.Vec3{-1, 0, 0},
		},
		{
			desc: "move forward after looking up",
			lookEvt: camera.LookEvent{
				Right: 0,
				Down:  -math.Pi / 2.0,
			},
			expected: mgl.Vec3{0, 1, 0},
		},
		{
			desc: "move forward after looking down",
			lookEvt: camera.LookEvent{
				Right: 0,
				Down:  math.Pi / 2.0,
			},
			expected: mgl.Vec3{0, -1, 0},
		},
		{
			desc: "move forward after looking right 180 degrees",
			lookEvt: camera.LookEvent{
				Right: math.Pi,
				Down:  0,
			},
			expected: mgl.Vec3{0, 0, 1},
		},
		{
			desc: "move forward after looking down 180 degrees",
			lookEvt: camera.LookEvent{
				Right: 0,
				Down:  math.Pi,
			},
			expected: mgl.Vec3{0, 0, 1},
		},
		{
			desc: "move forward after not looking anywhere else",
			lookEvt: camera.LookEvent{
				Right: 0,
				Down:  0,
			},
			expected: mgl.Vec3{0, 0, -1},
		},
		{
			desc: "move forward after doing 360 degree spin",
			lookEvt: camera.LookEvent{
				Right: 2 * math.Pi,
				Down:  0,
			},
			expected: mgl.Vec3{0, 0, -1},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			playerMod := &player.FnModule{}
			cameraMod := camera.New(playerMod, player.PositionEvent{
				X: 0,
				Y: 0,
				Z: 0,
			})
			var actual mgl.Vec3
			playerMod.FnUpdatePlayerPosition = func(posEvent player.PositionEvent) {
				actual = mgl.Vec3{
					posEvent.X,
					posEvent.Y,
					posEvent.Z,
				}
			}
			cameraMod.HandleLookEvent(tC.lookEvt)
			cameraMod.HandleMovementEvent(camera.MovementEvent{
				Direction: camera.MoveForwards,
				Pressed:   true,
			})
			cameraMod.Tick()
			if !withinErrorVec3(actual, tC.expected, errMargin) {
				t.Fatalf("expected new player position to be within %v of %v but was %v", errMargin, tC.expected, actual)
			}
		})
	}
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
		t.Run(fmt.Sprintf("updates player module direction when looking %+v", tC.looks), func(t *testing.T) {
			t.Parallel()

			var evt player.DirectionEvent
			playerMod := &player.FnModule{
				FnUpdatePlayerDirection: func(directionEvent player.DirectionEvent) {
					evt = directionEvent
				},
			}

			cameraMod := camera.New(playerMod, player.PositionEvent{})

			for _, look := range tC.looks {
				lookEvent := camera.LookEvent{
					Right: look.right,
					Down:  look.down,
				}
				cameraMod.HandleLookEvent(lookEvent)
			}

			expected := player.DirectionEvent{
				Rotation: tC.rotation,
			}
			if !withinErrorQuat(expected.Rotation, evt.Rotation, errMargin) {
				t.Fatalf("expected %v but got %v", expected, evt)
			}
		})
	}
}

func TestCameraInitialDirection(t *testing.T) {
	t.Parallel()
	expected := player.DirectionEvent{Rotation: mgl.QuatIdent()}
	var actual player.DirectionEvent
	playerMod := player.FnModule{
		FnUpdatePlayerDirection: func(dirEvent player.DirectionEvent) {
			actual = dirEvent
		},
	}
	camera.New(&playerMod, player.PositionEvent{X: 0, Y: 0, Z: 0})

	if actual != expected {
		t.Fatalf("expected quat %v but got %v", expected, actual)
	}
}

func TestCameraNilPlayer(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()

	camera.New(nil, player.PositionEvent{})
}
