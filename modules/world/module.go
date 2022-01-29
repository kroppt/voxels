package world

import "github.com/kroppt/voxels/modules/graphics"

type Module struct {
	c core
}

func New(graphicsMod graphics.Interface) *Module {
	return &Module{
		core{
			graphicsMod:  graphicsMod,
			chunksLoaded: map[ChunkEvent]struct{}{},
		},
	}
}
