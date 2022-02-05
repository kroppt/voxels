package view

import (
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
