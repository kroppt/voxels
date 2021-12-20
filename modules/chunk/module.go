package chunk

import "github.com/kroppt/voxels/modules/graphics"

// Module is a chunk module.
type Module struct {
	c core
}

// New creates a chunk module.
func New(graphicsMod graphics.Interface) *Module {
	core := core{
		graphicsMod: graphicsMod,
	}
	return &Module{
		core,
	}
}
