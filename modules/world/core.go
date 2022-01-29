package world

import "github.com/kroppt/voxels/modules/graphics"

type core struct {
	graphicsMod  graphics.Interface
	chunksLoaded map[ChunkEvent]struct{}
}

func (c *core) loadChunk(chunkEvent ChunkEvent) {
	c.chunksLoaded[chunkEvent] = struct{}{}
}

func (c *core) unloadChunk(chunkEvent ChunkEvent) {
	delete(c.chunksLoaded, chunkEvent)
}

func (c *core) countLoadedChunks() int {
	return len(c.chunksLoaded)
}

func (c *core) updateView(viewableChunks map[ChunkEvent]struct{}) {
	viewChunksForGraphics := map[graphics.ChunkEvent]struct{}{}
	for viewableChunk := range viewableChunks {
		viewChunksForGraphics[graphics.ChunkEvent(viewableChunk)] = struct{}{}
	}
	c.graphicsMod.UpdateViewableChunks(viewChunksForGraphics)
}
