package events_test

import (
	"fmt"
	"testing"

	"github.com/kroppt/voxels/modules/events"
	"github.com/kroppt/voxels/modules/player"
	"github.com/veandco/go-sdl2/sdl"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := events.New(nil, nil)
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

type fnGraphicsMod struct {
	FnPollEvent     func() (sdl.Event, bool)
	FnDestroyWindow func() error
}

func (fn fnGraphicsMod) PollEvent() (sdl.Event, bool) {
	return fn.FnPollEvent()
}

func (fn fnGraphicsMod) DestroyWindow() error {
	return fn.FnDestroyWindow()
}

type fnPlayerMod struct {
	FnHandleMovementEvent func(player.MovementEvent)
	FnHandleLookEvent     func(player.LookEvent)
}

func (fn fnPlayerMod) HandleMovementEvent(evt player.MovementEvent) {
	fn.FnHandleMovementEvent(evt)
}

func (fn fnPlayerMod) HandleLookEvent(evt player.LookEvent) {
	fn.FnHandleLookEvent(evt)
}

func TestModuleRouteEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns on quit event", func(t *testing.T) {
		t.Parallel()

		graphicsMod := fnGraphicsMod{
			FnPollEvent: func() (sdl.Event, bool) {
				return &sdl.QuitEvent{
					Type:      sdl.QUIT,
					Timestamp: 0,
				}, true
			},
			FnDestroyWindow: func() error {
				return nil
			},
		}
		mod := events.New(graphicsMod, nil)

		mod.RouteEvents()
	})

	t.Run("calls DestroyWindow on quit event", func(t *testing.T) {
		t.Parallel()

		destroyed := false
		graphicsMod := fnGraphicsMod{
			FnPollEvent: func() (sdl.Event, bool) {
				return &sdl.QuitEvent{
					Type:      sdl.QUIT,
					Timestamp: 0,
				}, true
			},
			FnDestroyWindow: func() error {
				destroyed = true
				return nil
			},
		}
		mod := events.New(graphicsMod, nil)

		mod.RouteEvents()

		if !destroyed {
			t.Fatal("expected tests to fail")
		}
	})

	testCases := []struct {
		message   string
		scancode  sdl.Scancode
		sym       sdl.Keycode
		direction player.MoveDirection
	}{
		{
			message:   "forward",
			scancode:  sdl.SCANCODE_W,
			sym:       sdl.K_w,
			direction: player.MoveForwards,
		},
		{
			message:   "backward",
			scancode:  sdl.SCANCODE_S,
			sym:       sdl.K_s,
			direction: player.MoveBackwards,
		},
		{
			message:   "left",
			scancode:  sdl.SCANCODE_A,
			sym:       sdl.K_a,
			direction: player.MoveLeft,
		},
		{
			message:   "right",
			scancode:  sdl.SCANCODE_D,
			sym:       sdl.K_d,
			direction: player.MoveRight,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run("creates a "+tC.message+" movement event from a keyboard event", func(t *testing.T) {
			t.Parallel()

			first := true
			moveKeyboardEvent := sdl.KeyboardEvent{
				Type:      sdl.KEYDOWN,
				Timestamp: 0,
				WindowID:  0,
				State:     sdl.PRESSED,
				Repeat:    0,
				Keysym: sdl.Keysym{
					Scancode: tC.scancode,
					Sym:      tC.sym,
					Mod:      sdl.KMOD_NONE,
				},
			}
			quitKeyboardEvent := sdl.QuitEvent{
				Type:      sdl.QUIT,
				Timestamp: 0,
			}
			graphicsMod := fnGraphicsMod{
				FnPollEvent: func() (sdl.Event, bool) {
					if first {
						first = false
						return &moveKeyboardEvent, true
					}
					return &quitKeyboardEvent, true
				},
				FnDestroyWindow: func() error {
					return nil
				},
			}

			movementEvent := player.MovementEvent{
				Direction: tC.direction,
			}
			var evtHandle *player.MovementEvent
			playerMod := &fnPlayerMod{
				FnHandleMovementEvent: func(evt player.MovementEvent) {
					evtHandle = &evt
				},
			}
			mod := events.New(graphicsMod, playerMod)

			mod.RouteEvents()

			if evtHandle == nil {
				t.Fatalf("expected %v but got %v", movementEvent, nil)
			}

			if movementEvent != *evtHandle {
				t.Fatalf("expected %v but got %v", movementEvent, *evtHandle)
			}
		})
	}

	testLookCases := []struct {
		xRel int32
		yRel int32
	}{
		{
			xRel: 1,
			yRel: 2,
		},
		{
			xRel: 0,
			yRel: 2,
		},
		{
			xRel: 1,
			yRel: 0,
		},
		{
			xRel: -1,
			yRel: 2,
		},
		{
			xRel: 1,
			yRel: -2,
		},
		{
			xRel: -1,
			yRel: -2,
		},
	}
	for _, tC := range testLookCases {
		tC := tC
		t.Run(fmt.Sprintf("convert sdl.MouseMotionEvent (%v, %v) to LookEvent", tC.xRel, tC.yRel), func(t *testing.T) {
			motionEvent := sdl.MouseMotionEvent{
				Type:      sdl.MOUSEMOTION,
				Timestamp: 0,
				WindowID:  0,
				Which:     0,
				State:     0,
				X:         0,
				Y:         0,
				XRel:      tC.xRel,
				YRel:      tC.yRel,
			}

			quitEvent := sdl.QuitEvent{
				Type:      sdl.QUIT,
				Timestamp: 0,
			}

			first := true
			graphicsMod := fnGraphicsMod{
				FnPollEvent: func() (sdl.Event, bool) {
					if first {
						first = false
						return &motionEvent, true
					}
					return &quitEvent, true
				},
				FnDestroyWindow: func() error {
					return nil
				},
			}
			var evtHandle *player.LookEvent
			playerMod := &fnPlayerMod{
				FnHandleLookEvent: func(evt player.LookEvent) {
					evtHandle = &evt
				},
			}
			mod := events.New(graphicsMod, playerMod)

			mod.RouteEvents()

			expectLookEvent := player.LookEvent{
				Right: tC.xRel,
				Down:  tC.yRel,
			}

			if evtHandle == nil {
				t.Fatalf("expected %v but got %v", expectLookEvent, nil)
			}

			if expectLookEvent != *evtHandle {
				t.Fatalf("expected %v but got %v", expectLookEvent, *evtHandle)
			}

		})
	}

}
