package world

import "github.com/kroppt/voxels/chunk"

type Generator interface {
	GenerateChunk(chunk.Position) chunk.Chunk
}

type FnGenerator struct {
	FnGenerateChunk func(chunk.Position) chunk.Chunk
}

func (fn *FnGenerator) GenerateChunk(pos chunk.Position) chunk.Chunk {
	if fn.FnGenerateChunk != nil {
		return fn.FnGenerateChunk(pos)
	}
	return chunk.New(pos, 1)
}
