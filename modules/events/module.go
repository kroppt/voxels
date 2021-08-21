package events

// Module is a synchronous event router.
type Module struct {
	c core
}

// New creates a synchronous events module.
func New(graphicsMod graphicsMod) *Module {
	return &Module{
		core{
			graphicsMod: graphicsMod,
			quit:        false,
		},
	}
}
