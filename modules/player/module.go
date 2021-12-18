package player

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

// Module is a player.
type Module struct {
	c core
}

type chunkMod interface {
	UpdatePosition(chunk.PositionEvent)
}

type graphicsMod interface {
	UpdatePlayerDirection(graphics.DirectionEvent)
}

// New creates a player.
func New(chunkMod chunkMod, graphicsMod graphicsMod) *Module {
	return &Module{
		core{
			chunkMod:    chunkMod,
			graphicsMod: graphicsMod,
			rot:         mgl.QuatIdent(),
		},
	}
}
