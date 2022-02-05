package world

import (
	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/chunk"
)

type Interface interface {
	LoadChunk(chunk.ChunkCoordinate)
	UnloadChunk(chunk.ChunkCoordinate)
	Quit()
	CountLoadedChunks() int
	GetBlockType(chunk.VoxelCoordinate) chunk.BlockType
	RemoveBlock(chunk.VoxelCoordinate)
	AddBlock(chunk.VoxelCoordinate, chunk.BlockType)
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

func (m *Module) Quit() {
	m.c.quit()
}

func (m *Module) CountLoadedChunks() int {
	return m.c.countLoadedChunks()
}

func (m *Module) GetBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	return m.c.getBlockType(pos)
}

func (m *Module) RemoveBlock(vc chunk.VoxelCoordinate) {
	m.c.removeBlock(vc)
}

func (m *Module) AddBlock(vc chunk.VoxelCoordinate, bt chunk.BlockType) {
	m.c.addBlock(vc, bt)
}

type FnModule struct {
	FnLoadChunk         func(chunk.ChunkCoordinate)
	FnUnloadChunk       func(chunk.ChunkCoordinate)
	FnQuit              func()
	FnCountLoadedChunks func() int
	FnGetBlockType      func(chunk.VoxelCoordinate) chunk.BlockType
	FnRemoveBlock       func(chunk.VoxelCoordinate)
	FnAddBlock          func(chunk.VoxelCoordinate, chunk.BlockType)
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

func (fn FnModule) GetBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	if fn.FnGetBlockType != nil {
		return fn.FnGetBlockType(pos)
	}
	return chunk.BlockTypeAir
}

func (fn FnModule) RemoveBlock(vc chunk.VoxelCoordinate) {
	if fn.FnRemoveBlock != nil {
		fn.FnRemoveBlock(vc)
	}
}

func (fn FnModule) AddBlock(vc chunk.VoxelCoordinate, bt chunk.BlockType) {
	if fn.FnAddBlock != nil {
		fn.FnAddBlock(vc, bt)
	}
}
