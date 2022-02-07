package world

import (
	"container/list"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/repositories/settings"
)

type Module struct {
	c core
}

func New(
	graphicsMod graphics.Interface,
	generator Generator,
	settingsRepo settings.Interface,
	cacheMod cache.Interface,
	viewMod view.Interface,
) *Module {
	if generator == nil {
		panic("world received a nil generator")
	}
	if settingsRepo == nil {
		panic("world received a nil settings repo")
	}
	if graphicsMod == nil {
		panic("world received a nil graphics module")
	}
	if viewMod == nil {
		panic("world received a nil view module")
	}
	return &Module{
		core{
			graphicsMod:    graphicsMod,
			generator:      generator,
			settingsRepo:   settingsRepo,
			cacheMod:       cacheMod,
			viewMod:        viewMod,
			loadedChunks:   map[chunk.ChunkCoordinate]*chunkState{},
			pendingActions: map[chunk.ChunkCoordinate]*list.List{},
		},
	}
}

type ParallelModule struct {
	do chan func()
	c  core
}

func NewParallel(
	graphicsMod graphics.Interface,
	generator Generator,
	settingsRepo settings.Interface,
	cacheMod cache.Interface,
	viewMod view.Interface,
) *ParallelModule {
	if generator == nil {
		panic("world received a nil generator")
	}
	if settingsRepo == nil {
		panic("world received a nil settings repo")
	}
	if graphicsMod == nil {
		panic("world received a nil graphics module")
	}
	if viewMod == nil {
		panic("world received a nil view module")
	}
	return &ParallelModule{
		do: make(chan func(), 1024),
		c: core{
			graphicsMod:    graphicsMod,
			generator:      generator,
			settingsRepo:   settingsRepo,
			cacheMod:       cacheMod,
			viewMod:        viewMod,
			loadedChunks:   map[chunk.ChunkCoordinate]*chunkState{},
			pendingActions: map[chunk.ChunkCoordinate]*list.List{},
		},
	}
}

func (m *ParallelModule) Run() {
	for f := range m.do {
		f()
	}
	m.Quit()
	m.c.viewMod.Close()
}

// Close stops the parallel execution.
//
// Close should be called when no more API calls will be used.
func (m *ParallelModule) Close() {
	close(m.do)
}
