package input

// Module is a synchronous input router.
type Module struct {
	c core
}

// New creates a synchronous input module.
func New(graphicsMod graphicsMod, playerMod playerMod) *Module {
	return &Module{
		core{
			graphicsMod: graphicsMod,
			playerMod:   playerMod,
			quit:        false,
		},
	}
}
