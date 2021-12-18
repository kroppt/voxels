package input

import (
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
)

// Module is a synchronous input router.
type Module struct {
	c core
}

type playerMod interface {
	HandleMovementEvent(player.MovementEvent)
	HandleLookEvent(player.LookEvent)
}

type settingsRepo interface {
	GetFOV() float64
	GetResolution() (uint32, uint32)
}

// New creates a synchronous input module.
func New(
	graphicsMod graphics.Interface,
	playerMod playerMod,
	settingsRepo settingsRepo,
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
