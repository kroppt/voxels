package chunk

// Module is a chunk module.
type Module struct {
	c core
}

// New creates a chunk module.
func New(graphicsMod graphicsMod) *Module {
	core := core{
		graphicsMod: graphicsMod,
	}
	core.init()
	return &Module{
		core,
	}
}
