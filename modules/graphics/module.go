package graphics

// Module is a synchronous graphics renderer.
type Module struct {
	c core
}

// New creates a synchronous events module.
func New() *Module {
	return &Module{
		core{
			window: nil,
		},
	}
}
