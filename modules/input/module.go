package input

import (
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
	playerMod player.Interface,
	settingsRepo settings.Interface,
) *Module {
	return &Module{
		core{
			graphicsMod:  graphicsMod,
			playerMod:    playerMod,
			settingsRepo: settingsRepo,
			quit:         false,
		},
	}
}
