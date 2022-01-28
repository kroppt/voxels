package world

import "github.com/kroppt/voxels/modules/graphics"

type core struct {
	graphicsMod graphics.Interface
}

func (c *core) updateView(viewableChunks map[ChunkEvent]struct{}) {
	viewChunksForGraphics := map[graphics.ChunkEvent]struct{}{}
	for viewableChunk := range viewableChunks {
		viewChunksForGraphics[graphics.ChunkEvent(viewableChunk)] = struct{}{}
	}
	c.graphicsMod.UpdateViewableChunks(viewChunksForGraphics)
}
