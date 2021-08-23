package input

// Module is a synchronous input router.
type Module struct {
	c core
}

// New creates a synchronous input module.
func New(
	graphicsMod graphicsMod,
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
