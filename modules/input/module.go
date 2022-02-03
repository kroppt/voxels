package input

import (
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a synchronous input router.
type Module struct {
	c core
}

// New creates a synchronous input module.
func New(
	graphicsMod graphics.Interface,
	cameraMod camera.Interface,
	settingsRepo settings.Interface,
	playerMod player.Interface,
) *Module {
	return &Module{
		core{
			graphicsMod:  graphicsMod,
			cameraMod:    cameraMod,
			settingsRepo: settingsRepo,
			playerMod:    playerMod,
			quit:         false,
		},
	}
}
