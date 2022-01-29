package world

import "github.com/kroppt/voxels/chunk"

type Interface interface {
	LoadChunk(chunk.Position)
	UnloadChunk(chunk.Position)
	UpdateView(map[chunk.Position]struct{})
	CountLoadedChunks() int
}

func (m *Module) LoadChunk(pos chunk.Position) {
	m.c.loadChunk(pos)
}

func (m *Module) UnloadChunk(pos chunk.Position) {
	m.c.unloadChunk(pos)
}

func (m *Module) UpdateView(viewableChunks map[chunk.Position]struct{}) {
	m.c.updateView(viewableChunks)
}

func (m *Module) CountLoadedChunks() int {
	return m.c.countLoadedChunks()
}

type FnModule struct {
	FnLoadChunk         func(chunk.Position)
	FnUnloadChunk       func(chunk.Position)
	FnUpdateView        func(map[chunk.Position]struct{})
	FnCountLoadedChunks func() int
}

func (fn *FnModule) LoadChunk(pos chunk.Position) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(pos)
	}
}
func (fn *FnModule) UnloadChunk(pos chunk.Position) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(pos)
	}
}
func (fn *FnModule) UpdateView(viewChunks map[chunk.Position]struct{}) {
	if fn.FnUpdateView != nil {
		fn.FnUpdateView(viewChunks)
	}
}
func (fn *FnModule) CountLoadedChunks() int {
	if fn.FnCountLoadedChunks != nil {
		return fn.FnCountLoadedChunks()
	}
	return 0
}
