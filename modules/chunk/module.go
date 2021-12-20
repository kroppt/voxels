package chunk

import (
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a chunk module.
type Module struct {
	c core
}

// New creates a chunk module.
func New(graphicsMod graphics.Interface, settingsMod settings.Interface) *Module {
	if settingsMod == nil {
		panic("settings is required to be non-nil")
	}
	core := core{
		graphicsMod: graphicsMod,
		settingsMod: settingsMod,
	}
	core.init()
	return &Module{
		core,
	}
}
