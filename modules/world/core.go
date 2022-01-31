package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod  graphics.Interface
	generator    Generator
	settingsRepo settings.Interface
	cacheMod     cache.Interface
	chunksLoaded map[chunk.Position]chunk.Chunk
}

func (c *core) loadChunk(pos chunk.Position) {
	if _, ok := c.chunksLoaded[pos]; ok {
		panic("tried to load already-loaded chunk")
	}
	ch, ok := c.cacheMod.Load(pos)
	if !ok {
		ch = c.generator.GenerateChunk(pos)
	}
	c.chunksLoaded[pos] = ch
	c.graphicsMod.LoadChunk(ch)
}

func (c *core) unloadChunk(pos chunk.Position) {
	if _, ok := c.chunksLoaded[pos]; !ok {
		panic("tried to unload a chunk that is not loaded")
	}
	c.cacheMod.Save(c.chunksLoaded[pos])
	delete(c.chunksLoaded, pos)
	c.graphicsMod.UnloadChunk(pos)
}

func (c *core) quit() {
	for key := range c.chunksLoaded {
		c.unloadChunk(key)
	}
	c.cacheMod.Close()
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
