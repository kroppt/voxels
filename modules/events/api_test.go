package events_test

import (
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
}

func (fn fnPlayerMod) HandleMovementEvent(evt player.MovementEvent) {
	fn.FnHandleMovementEvent(evt)
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

	for _, tc := range testCases {
		t.Run("creates a "+tc.message+" movement event from a keyboard event", func(t *testing.T) {
			t.Parallel()

			first := true
			moveKeyboardEvent := sdl.KeyboardEvent{
				Type:      sdl.KEYDOWN,
				Timestamp: 0,
				WindowID:  0,
				State:     sdl.PRESSED,
				Repeat:    0,
				Keysym: sdl.Keysym{
					Scancode: tc.scancode,
					Sym:      tc.sym,
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
				Direction: tc.direction,
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

}
