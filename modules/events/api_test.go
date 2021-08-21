package events_test

import (
	"testing"

	"github.com/kroppt/voxels/modules/events"
	"github.com/veandco/go-sdl2/sdl"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := events.New(nil)
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
		mod := events.New(graphicsMod)

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
		mod := events.New(graphicsMod)

		mod.RouteEvents()

		if !destroyed {
			t.Fatal("expected tests to fail")
		}
	})

}
