package player

import (
	"context"

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
		panic("view is required to be non-nil")
	}
	if worldMod == nil {
		panic("world is required to be non-nil")
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

type ParallelModule struct {
	do chan func()
	c  core
}

// NewParallel creates a parallel player module.
func NewParallel(worldMod world.Interface, settingsMod settings.Interface, viewMod view.Interface) *ParallelModule {
	if settingsMod == nil {
		panic("settings is required to be non-nil")
	}
	if viewMod == nil {
		panic("view is required to be non-nil")
	}
	if worldMod == nil {
		panic("world is required to be non-nil")
	}
	core := core{
		worldMod:    worldMod,
		settingsMod: settingsMod,
		viewMod:     viewMod,
		firstLoad:   true,
	}
	return &ParallelModule{
		do: make(chan func()),
		c:  core,
	}
}

func (m *ParallelModule) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-m.do:
			f()
		}
	}
}
