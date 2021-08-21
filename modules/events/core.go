package events

import (
	"github.com/kroppt/voxels/log"
	"github.com/veandco/go-sdl2/sdl"
)

type graphicsMod interface {
	DestroyWindow() error
	PollEvent() (sdl.Event, bool)
}

type core struct {
	graphicsMod graphicsMod
	quit        bool
}

func (m *core) routeEvents() {
	for !m.quit {
		evt, ok := m.graphicsMod.PollEvent()
		if !ok {
			continue
		}

		m.routeEvent(evt)
	}
}

func (m *core) routeEvent(e sdl.Event) {
	switch e.(type) {
	case *sdl.QuitEvent:
		err := m.graphicsMod.DestroyWindow()
		if err != nil {
			log.Warnf("error closing window: %v", err)
		}
		m.quit = true
	}
}
