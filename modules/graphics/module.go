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
			loadedChunks:   map[chunk.ChunkCoordinate]*glObject{},
			viewableChunks: map[chunk.ChunkCoordinate]struct{}{},
		},
	}
}

type ParallelModule struct {
	do chan func()
	c  core
}

func NewParallel(settingsRepo settings.Interface) *ParallelModule {
	return &ParallelModule{
		do: make(chan func(), 10),
		c: core{
			window:         nil,
			settingsRepo:   settingsRepo,
			loadedChunks:   map[chunk.ChunkCoordinate]*glObject{},
			viewableChunks: map[chunk.ChunkCoordinate]struct{}{},
		},
	}
}

func (m *ParallelModule) Run() {
	for f := range m.do {
		f()
	}
	m.DestroyWindow()
}

// Close stops the parallel execution.
//
// Close should be called when no more API calls will be used.
func (m *ParallelModule) Close() {
	close(m.do)
}
