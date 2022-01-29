package world

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type core struct {
	graphicsMod  graphics.Interface
	generator    Generator
	chunksLoaded map[chunk.Position]struct{}
}

func (c *core) loadChunk(pos chunk.Position) {
	c.chunksLoaded[pos] = struct{}{}
	c.graphicsMod.LoadChunk(c.generator.GenerateChunk(pos))
}

func (c *core) unloadChunk(pos chunk.Position) {
	delete(c.chunksLoaded, pos)
	c.graphicsMod.UnloadChunk(pos)
}

func (c *core) countLoadedChunks() int {
	return len(c.chunksLoaded)
}
