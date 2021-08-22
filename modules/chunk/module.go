package chunk

// Module is a chunk module.
type Module struct {
	c core
}

// New creates a chunk module.
func New() *Module {
	return &Module{
		core{},
	}
}
