package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type Module struct {
	c core
}

func New(graphicsMod graphics.Interface, generator Generator) *Module {
	return &Module{
		core{
			graphicsMod:  graphicsMod,
			generator:    generator,
			chunksLoaded: map[chunk.Position]struct{}{},
		},
	}
}
