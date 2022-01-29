package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/player"
)

// Module is a camera.
type Module struct {
	c core
}

// New creates a camera.
func New(playerMod player.Interface, initialPos player.PositionEvent) *Module {
	if playerMod != nil {
		playerMod.UpdatePlayerPosition(initialPos)
	}
	return &Module{
		core{
			playerMod: playerMod,
			rot:       mgl.QuatIdent(),
		},
	}
}
