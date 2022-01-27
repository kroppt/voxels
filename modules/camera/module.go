package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

// Module is a camera.
type Module struct {
	c core
}

// New creates a camera.
func New(chunkMod chunk.Interface, graphicsMod graphics.Interface) *Module {
	return &Module{
		core{
			chunkMod:    chunkMod,
			graphicsMod: graphicsMod,
			rot:         mgl.QuatIdent(),
		},
	}
}
