package world

type Interface interface {
	LoadChunk(ChunkEvent)
	UnloadChunk(ChunkEvent)
	UpdateView(map[ChunkEvent]struct{})
	CountLoadedChunks() int
}

type ChunkEvent struct {
	PositionX int32
	PositionY int32
	PositionZ int32
}

func (m *Module) LoadChunk(chunkEvent ChunkEvent) {
	m.c.loadChunk(chunkEvent)
}

func (m *Module) UnloadChunk(chunkEvent ChunkEvent) {
	m.c.unloadChunk(chunkEvent)
}

func (m *Module) UpdateView(viewableChunks map[ChunkEvent]struct{}) {
	m.c.updateView(viewableChunks)
}

func (m *Module) CountLoadedChunks() int {
	return m.c.countLoadedChunks()
}

type FnModule struct {
	FnLoadChunk         func(ChunkEvent)
	FnUnloadChunk       func(ChunkEvent)
	FnUpdateView        func(map[ChunkEvent]struct{})
	FnCountLoadedChunks func() int
}

func (fn *FnModule) LoadChunk(evt ChunkEvent) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(evt)
	}
}
func (fn *FnModule) UnloadChunk(evt ChunkEvent) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(evt)
	}
}
func (fn *FnModule) UpdateView(viewChunks map[ChunkEvent]struct{}) {
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
