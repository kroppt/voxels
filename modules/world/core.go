package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod  graphics.Interface
	generator    Generator
	settingsRepo settings.Interface
	chunksLoaded map[chunk.Position]chunk.Chunk
}

func (c *core) loadChunk(pos chunk.Position) {
	if _, ok := c.chunksLoaded[pos]; ok {
		panic("tried to load already-loaded chunk")
	}
	generatedChunk := c.generator.GenerateChunk(pos)
	c.chunksLoaded[pos] = generatedChunk
	c.graphicsMod.LoadChunk(generatedChunk)
}

func (c *core) unloadChunk(pos chunk.Position) {
	if _, ok := c.chunksLoaded[pos]; !ok {
		panic("tried to unload a chunk that is not loaded")
	}
	delete(c.chunksLoaded, pos)
	c.graphicsMod.UnloadChunk(pos)
}

func (c *core) countLoadedChunks() int {
	return len(c.chunksLoaded)
}

func (c *core) setBlockType(pos chunk.VoxelCoordinate, btype chunk.BlockType) {
	key := chunk.VoxelCoordToChunkCoord(pos, c.settingsRepo.GetChunkSize())
	if _, ok := c.chunksLoaded[key]; !ok {
		panic("tried to set block in non-loaded chunk")
	}
	c.chunksLoaded[key].SetBlockType(pos, btype)
}

func (c *core) getBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	key := chunk.VoxelCoordToChunkCoord(pos, c.settingsRepo.GetChunkSize())
	if _, ok := c.chunksLoaded[key]; !ok {
		panic("tried to get block from non-loaded chunk")
	}
	return c.chunksLoaded[key].BlockType(pos)
}
