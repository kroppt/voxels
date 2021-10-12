package chunk

// Module is a chunk module.
type Module struct {
	c core
}

// New creates a chunk module.
func New(graphicsMod graphicsMod, chunkSize uint) *Module {
	if chunkSize == 0 {
		panic("chunk size cannot be 0")
	}

	core := core{
		graphicsMod: graphicsMod,
	}
	core.init(chunkSize)
	return &Module{
		core,
	}
}
