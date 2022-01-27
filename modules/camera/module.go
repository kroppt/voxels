package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
)

// Module is a camera.
type Module struct {
	c core
}

// New creates a camera.
func New(playerMod player.Interface, graphicsMod graphics.Interface) *Module {
	return &Module{
		core{
			playerMod:   playerMod,
			graphicsMod: graphicsMod,
			rot:         mgl.QuatIdent(),
		},
	}
}
