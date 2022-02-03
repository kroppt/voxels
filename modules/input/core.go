package input

import (
	"math"

	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/veandco/go-sdl2/sdl"
)

type core struct {
	graphicsMod  graphics.Interface
	cameraMod    camera.Interface
	settingsRepo settings.Interface
	playerMod    player.Interface
	quit         bool
}

func (m *core) routeEvents() bool {
	for !m.quit {
		evt, ok := m.graphicsMod.PollEvent()
		if !ok {
			return true
		}

		m.routeEvent(evt)
	}
	return false
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
		pressed := evt.Type == sdl.KEYDOWN
		switch evt.Keysym.Scancode {
		case sdl.SCANCODE_W:
			forward := camera.MovementEvent{
				Direction: camera.MoveForwards,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(forward)
		case sdl.SCANCODE_D:
			right := camera.MovementEvent{
				Direction: camera.MoveRight,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(right)
		case sdl.SCANCODE_S:
			back := camera.MovementEvent{
				Direction: camera.MoveBackwards,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(back)
		case sdl.SCANCODE_A:
			left := camera.MovementEvent{
				Direction: camera.MoveLeft,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(left)
		case sdl.SCANCODE_SPACE:
			up := camera.MovementEvent{
				Direction: camera.MoveUp,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(up)
		case sdl.SCANCODE_LSHIFT:
			down := camera.MovementEvent{
				Direction: camera.MoveDown,
				Pressed:   pressed,
			}
			m.cameraMod.HandleMovementEvent(down)
		}

	case *sdl.MouseMotionEvent:
		xRad, yRad := m.pixelsToRadians(evt.XRel, evt.YRel)
		lookEvt := camera.LookEvent{
			Right: xRad,
			Down:  yRad,
		}
		if evt.State == sdl.BUTTON_LEFT {
			m.cameraMod.HandleLookEvent(lookEvt)
		}
	case *sdl.MouseWheelEvent:
		if evt.Y < 0 {
			m.playerMod.UpdatePlayerAction(player.ActionEvent{Scroll: player.ScrollDown})
		} else {
			m.playerMod.UpdatePlayerAction(player.ActionEvent{Scroll: player.ScrollUp})
		}
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
