package player

import (
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a player module.
type Module struct {
	c core
}

// New creates a player module.
func New(worldMod world.Interface, settingsMod settings.Interface, graphicsMod graphics.Interface) *Module {
	if settingsMod == nil {
		panic("settings is required to be non-nil")
	}
	core := core{
		worldMod:    worldMod,
		settingsMod: settingsMod,
		graphicsMod: graphicsMod,
		firstLoad:   true,
	}
	return &Module{
		core,
	}
}
