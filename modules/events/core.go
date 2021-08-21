package events

import (
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/player"
	"github.com/veandco/go-sdl2/sdl"
)

type graphicsMod interface {
	DestroyWindow() error
	PollEvent() (sdl.Event, bool)
}

type playerMod interface {
	HandleMovementEvent(player.MovementEvent)
}

type core struct {
	graphicsMod graphicsMod
	playerMod   playerMod
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
	switch evt := e.(type) {
	case *sdl.QuitEvent:
		err := m.graphicsMod.DestroyWindow()
		if err != nil {
			log.Warnf("error closing window: %v", err)
		}
		m.quit = true
	case *sdl.KeyboardEvent:
		if evt.Type != sdl.KEYDOWN {
			break
		}
		switch evt.Keysym.Scancode {
		case sdl.SCANCODE_W:
			forward := player.MovementEvent{
				Direction: player.MoveForwards,
			}
			m.playerMod.HandleMovementEvent(forward)
		case sdl.SCANCODE_D:
			right := player.MovementEvent{
				Direction: player.MoveRight,
			}
			m.playerMod.HandleMovementEvent(right)
		case sdl.SCANCODE_S:
			back := player.MovementEvent{
				Direction: player.MoveBackwards,
			}
			m.playerMod.HandleMovementEvent(back)
		case sdl.SCANCODE_A:
			left := player.MovementEvent{
				Direction: player.MoveLeft,
			}
			m.playerMod.HandleMovementEvent(left)
		}
	}
}
