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

func (c *core) voxelPosToChunkPos(pos chunk.VoxelCoordinate) chunk.Position {
	x, y, z := pos.X, pos.Y, pos.Z
	chunkSize := int32(c.settingsRepo.GetChunkSize())
	if pos.X < 0 {
		x++
	}
	if pos.Y < 0 {
		y++
	}
	if pos.Z < 0 {
		z++
	}
	x /= chunkSize
	y /= chunkSize
	z /= chunkSize
	if pos.X < 0 {
		x--
	}
	if pos.Y < 0 {
		y--
	}
	if pos.Z < 0 {
		z--
	}
	return chunk.Position{X: x, Y: y, Z: z}
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
	key := c.voxelPosToChunkPos(pos)
	if _, ok := c.chunksLoaded[key]; !ok {
		panic("tried to set block in non-loaded chunk")
	}
	c.chunksLoaded[key].SetBlockType(pos, btype)
}

func (c *core) getBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	key := c.voxelPosToChunkPos(pos)
	if _, ok := c.chunksLoaded[key]; !ok {
		panic("tried to get block from non-loaded chunk")
	}
	return c.chunksLoaded[key].BlockType(pos)
}
