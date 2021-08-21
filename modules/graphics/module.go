package graphics

import "github.com/veandco/go-sdl2/sdl"

// Module is a synchronous graphics renderer.
type Module struct {
	c core
}

// New creates a synchronous events module.
func New() *Module {
	return &Module{
		core{
			window: &sdl.Window{},
		},
	}
}
