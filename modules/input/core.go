package input

import (
	"math"

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
	HandleLookEvent(player.LookEvent)
}

type settingsRepo interface {
	GetFOV() float64
	GetResolution() (uint32, uint32)
}

type core struct {
	graphicsMod  graphicsMod
	playerMod    playerMod
	settingsRepo settingsRepo
	quit         bool
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
	case *sdl.MouseMotionEvent:
		xRad, yRad := m.pixelsToRadians(evt.XRel, evt.YRel)
		lookEvt := player.LookEvent{
			Right: xRad,
			Down:  yRad,
		}
		m.playerMod.HandleLookEvent(lookEvt)
	}

}

func (m *core) pixelsToRadians(xRel, yRel int32) (float64, float64) {
	const nearDistance = 0.1
	fovY := m.settingsRepo.GetFOV() * math.Pi / 180
	_, screenHeight := m.settingsRepo.GetResolution()
	nearHeight := 2 * nearDistance * math.Tan(fovY/2)
	radPerPixel := math.Atan(nearHeight / float64(screenHeight) / 0.1)
	return radPerPixel * float64(xRel), radPerPixel * float64(yRel)
}
