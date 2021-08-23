package player

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
		},
	}
}
