package player

// Module is a player.
type Module struct {
	c core
}

// New creates a player.
func New() *Module {
	return &Module{
		core{},
	}
}
