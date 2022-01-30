package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type Module struct {
	c core
}

func New(graphicsMod graphics.Interface, generator Generator, settingsRepo settings.Interface) *Module {
	if generator == nil {
		panic("world received a nil generator")
	}
	if settingsRepo == nil {
		panic("world received a nil settings repo")
	}
	if graphicsMod == nil {
		panic("world received a nil graphics module")
	}
	return &Module{
		core{
			graphicsMod:  graphicsMod,
			generator:    generator,
			settingsRepo: settingsRepo,
			chunksLoaded: map[chunk.Position]chunk.Chunk{},
		},
	}
}
