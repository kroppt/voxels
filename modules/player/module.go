package player

import (
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a player module.
type Module struct {
	c core
}

// New creates a player module.
func New(worldMod world.Interface, settingsMod settings.Interface, chunkSize uint32) *Module {
	if settingsMod == nil {
		panic("settings is required to be non-nil")
	}
	core := core{
		worldMod:    worldMod,
		settingsMod: settingsMod,
		chunkSize:   chunkSize,
	}
	core.init()
	return &Module{
		core,
	}
}
