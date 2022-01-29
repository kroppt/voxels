package world

import "github.com/kroppt/voxels/chunk"

type Generator interface {
	GenerateChunk(ChunkEvent) chunk.Chunk
}

type FnGenerator struct {
	FnGenerateChunk func(ChunkEvent) chunk.Chunk
}

func (fn *FnGenerator) GenerateChunk(chunkEvent ChunkEvent) chunk.Chunk {
	if fn.FnGenerateChunk != nil {
		return fn.FnGenerateChunk(chunkEvent)
	}
	return chunk.New(chunk.Position{
		X: chunkEvent.PositionX,
		Y: chunkEvent.PositionY,
		Z: chunkEvent.PositionZ,
	}, 0)
}
