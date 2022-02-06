package view

import (
	"context"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type Module struct {
	c *core
}

func New(graphicsMod graphics.Interface, settingsRepo settings.Interface) *Module {
	if graphicsMod == nil {
		panic("view module received nil graphics module")
	}
	if settingsRepo == nil {
		panic("view module received nil settings repo")
	}
	return &Module{
		c: &core{
			graphicsMod:  graphicsMod,
			settingsRepo: settingsRepo,
			trees:        map[chunk.ChunkCoordinate]*Octree{},
		},
	}
}

type ParallelModule struct {
	do chan func()
	c  *core
}

func NewParallel(graphicsMod graphics.Interface, settingsRepo settings.Interface) *ParallelModule {
	if graphicsMod == nil {
		panic("view module received nil graphics module")
	}
	if settingsRepo == nil {
		panic("view module received nil settings repo")
	}
	return &ParallelModule{
		do: make(chan func()),
		c: &core{
			graphicsMod:  graphicsMod,
			settingsRepo: settingsRepo,
			trees:        map[chunk.ChunkCoordinate]*Octree{},
		},
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
