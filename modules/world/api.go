package world

import "github.com/kroppt/voxels/chunk"

type Interface interface {
	LoadChunk(chunk.Position)
	UnloadChunk(chunk.Position)
	UnloadAllChunks()
	CountLoadedChunks() int
	SetBlockType(chunk.VoxelCoordinate, chunk.BlockType)
	GetBlockType(chunk.VoxelCoordinate) chunk.BlockType
}

func (m *Module) LoadChunk(pos chunk.Position) {
	m.c.loadChunk(pos)
}

func (m *Module) UnloadChunk(pos chunk.Position) {
	m.c.unloadChunk(pos)
}

func (m *Module) UnloadAllChunks() {
	m.c.unloadAllChunks()
}

func (m *Module) CountLoadedChunks() int {
	return m.c.countLoadedChunks()
}

func (m *Module) SetBlockType(pos chunk.VoxelCoordinate, btype chunk.BlockType) {
	m.c.setBlockType(pos, btype)
}

func (m *Module) GetBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	return m.c.getBlockType(pos)
}

type FnModule struct {
	FnLoadChunk         func(chunk.Position)
	FnUnloadChunk       func(chunk.Position)
	FnUnloadAllChunks   func()
	FnCountLoadedChunks func() int
	FnSetBlockType      func(chunk.VoxelCoordinate, chunk.BlockType)
	FnGetBlockType      func(chunk.VoxelCoordinate) chunk.BlockType
}

func (fn FnModule) LoadChunk(pos chunk.Position) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(pos)
	}
}

func (fn FnModule) UnloadChunk(pos chunk.Position) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(pos)
	}
}

func (fn FnModule) UnloadAllChunks() {
	if fn.FnUnloadAllChunks != nil {
		fn.FnUnloadAllChunks()
	}
}

func (fn FnModule) CountLoadedChunks() int {
	if fn.FnCountLoadedChunks != nil {
		return fn.FnCountLoadedChunks()
	}
	return 0
}

func (fn FnModule) SetBlockType(pos chunk.VoxelCoordinate, btype chunk.BlockType) {
	if fn.FnSetBlockType != nil {
		fn.FnSetBlockType(pos, btype)
	}
}

func (fn FnModule) GetBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	if fn.FnGetBlockType != nil {
		return fn.FnGetBlockType(pos)
	}
	return chunk.BlockTypeAir
}
