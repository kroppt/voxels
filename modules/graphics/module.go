package graphics

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/repositories/settings"
)

// Module is a synchronous graphics renderer.
type Module struct {
	c core
}

// New creates a synchronous events module.
func New(settingsRepo settings.Interface) *Module {
	return &Module{
		core{
			window:         nil,
			settingsRepo:   settingsRepo,
			loadedChunks:   map[chunk.ChunkCoordinate]*ChunkObject{},
			viewableChunks: map[chunk.ChunkCoordinate]struct{}{},
		},
	}
}
