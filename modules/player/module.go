package player

import (
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a player module.
type Module struct {
	c core
}

// New creates a player module.
func New(worldMod world.Interface, settingsMod settings.Interface, viewMod view.Interface) *Module {
	if settingsMod == nil {
		panic("settings is required to be non-nil")
	}
	if viewMod == nil {
		panic("player module received a nil view module")
	}
	if worldMod == nil {
		panic("player module received a nil world module")
	}
	core := core{
		worldMod:    worldMod,
		settingsMod: settingsMod,
		viewMod:     viewMod,
		firstLoad:   true,
	}
	return &Module{
		core,
	}
}
