package world

import (
	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/chunk"
)

type Interface interface {
	LoadChunk(chunk.ChunkCoordinate)
	UnloadChunk(chunk.ChunkCoordinate)
	UpdateView(ViewState)
	Quit()
	CountLoadedChunks() int
	SetBlockType(chunk.VoxelCoordinate, chunk.BlockType)
	GetBlockType(chunk.VoxelCoordinate) chunk.BlockType
}

type ViewState struct {
	Pos mgl.Vec3
	Dir mgl.Quat
}

func (m *Module) LoadChunk(pos chunk.ChunkCoordinate) {
	m.c.loadChunk(pos)
}

func (m *Module) UnloadChunk(pos chunk.ChunkCoordinate) {
	m.c.unloadChunk(pos)
}

func (m *Module) UpdateView(viewState ViewState) {
	m.c.updateView(viewState)
}

func (m *Module) Quit() {
	m.c.quit()
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
	FnLoadChunk         func(chunk.ChunkCoordinate)
	FnUnloadChunk       func(chunk.ChunkCoordinate)
	FnUpdateView        func(ViewState)
	FnQuit              func()
	FnCountLoadedChunks func() int
	FnSetBlockType      func(chunk.VoxelCoordinate, chunk.BlockType)
	FnGetBlockType      func(chunk.VoxelCoordinate) chunk.BlockType
}

func (fn FnModule) LoadChunk(pos chunk.ChunkCoordinate) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(pos)
	}
}

func (fn FnModule) UnloadChunk(pos chunk.ChunkCoordinate) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(pos)
	}
}

func (fn FnModule) UpdateView(viewState ViewState) {
	if fn.FnUpdateView != nil {
		fn.FnUpdateView(viewState)
	}
}

func (fn FnModule) Quit() {
	if fn.FnQuit != nil {
		fn.FnQuit()
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
