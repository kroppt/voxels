package player

import mgl "github.com/go-gl/mathgl/mgl64"

// Module is a player.
type Module struct {
	c core
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
