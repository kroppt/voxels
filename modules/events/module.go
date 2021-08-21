package events

// Module is a synchronous event router.
type Module struct {
	c core
}

// New creates a synchronous events module.
func New(graphicsMod graphicsMod, playerMod playerMod) *Module {
	return &Module{
		core{
			graphicsMod: graphicsMod,
			playerMod:   playerMod,
			quit:        false,
		},
	}
}
