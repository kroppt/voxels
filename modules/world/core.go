package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type core struct {
	graphicsMod  graphics.Interface
	generator    Generator
	chunksLoaded map[ChunkEvent]struct{}
}

func (c *core) loadChunk(chunkEvent ChunkEvent) {
	c.chunksLoaded[chunkEvent] = struct{}{}
	c.graphicsMod.LoadChunk(c.generator.GenerateChunk(chunkEvent))
}

func (c *core) unloadChunk(chunkEvent ChunkEvent) {
	delete(c.chunksLoaded, chunkEvent)
	c.graphicsMod.UnloadChunk(chunk.Position{
		X: chunkEvent.PositionX,
		Y: chunkEvent.PositionY,
		Z: chunkEvent.PositionZ,
	})
}

func (c *core) countLoadedChunks() int {
	return len(c.chunksLoaded)
}

func (c *core) updateView(viewableChunks map[ChunkEvent]struct{}) {
	viewChunksForGraphics := map[chunk.Position]struct{}{}
	for viewableChunk := range viewableChunks {
		viewChunksForGraphics[chunk.Position{
			X: viewableChunk.PositionX,
			Y: viewableChunk.PositionY,
			Z: viewableChunk.PositionZ,
		}] = struct{}{}
	}
	c.graphicsMod.UpdateViewableChunks(viewChunksForGraphics)
}
