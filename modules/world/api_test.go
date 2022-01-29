package world_test

import (
	"testing"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
)

func TestWorldLoadedChunkCount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc          string
		loadChunks    []chunk.Position
		unloadChunks  []chunk.Position
		expectedCount int
	}{
		{
			desc:          "world starts with no loaded chunks",
			expectedCount: 0,
		},
		{
			desc: "world should load one chunk",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should load and unload the same chunk",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			unloadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should not load the same chunk twice",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should not unload the same chunk twice",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			unloadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should load two different chunks",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
			},
			expectedCount: 2,
		},
		{
			desc: "world cannot unload a chunk if it has none",
			unloadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 0,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			graphicsMod := graphics.FnModule{}
			worldMod := world.New(graphicsMod, &world.FnGenerator{})
			for _, loadChunk := range tC.loadChunks {
				worldMod.LoadChunk(loadChunk)
			}
			for _, unloadChunk := range tC.unloadChunks {
				worldMod.UnloadChunk(unloadChunk)
			}
			actual := worldMod.CountLoadedChunks()
			if actual != tC.expectedCount {
				t.Fatalf("expected %v chunks to be loaded but got %v", tC.expectedCount, actual)
			}
		})
	}
}

func TestWorldLoadChunkPassesToGraphics(t *testing.T) {
	t.Parallel()
	var actual chunk.Position
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			actual = ch.Position()
		},
	}

	worldMod := world.New(graphicsMod, &world.FnGenerator{})
	worldMod.LoadChunk(chunk.Position{X: 1, Y: 2, Z: 3})
	expected := chunk.Position{X: 1, Y: 2, Z: 3}
	if actual != expected {
		t.Fatalf("expected graphics to receive %v but got %v", expected, actual)
	}
}

func TestWorldUnloadChunkPassesToGraphics(t *testing.T) {
	t.Parallel()
	var actual chunk.Position
	graphicsMod := graphics.FnModule{
		FnUnloadChunk: func(pos chunk.Position) {
			actual = pos
		},
	}
	worldMod := world.New(graphicsMod, nil)
	worldMod.UnloadChunk(chunk.Position{X: 1, Y: 2, Z: 3})
	expected := chunk.Position{X: 1, Y: 2, Z: 3}
	if actual != expected {
		t.Fatalf("expected graphics to receive %v but got %v", expected, actual)
	}
}

func TestWorldGeneration(t *testing.T) {
	t.Parallel()
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(_ chunk.Position) chunk.Chunk {
			newChunk := chunk.New(chunk.Position{X: 0, Y: 0, Z: 0}, 1)
			newChunk.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeDirt)
			return newChunk
		},
	}
	expected := chunk.BlockTypeDirt
	var actual chunk.BlockType
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			actual = ch.BlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
		},
	}
	worldMod := world.New(graphicsMod, testGen)
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})

	if actual != expected {
		t.Fatalf("expected to retrieve block type %v but got %v", expected, actual)
	}

}
