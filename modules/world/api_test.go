package world_test

import (
	"reflect"
	"testing"

	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
)

func TestWorldUpdateChunkView(t *testing.T) {
	t.Parallel()
	expectedViewableChunks := map[graphics.ChunkEvent]struct{}{
		{1, 2, 3}: {},
	}
	var actualViewableChunks map[graphics.ChunkEvent]struct{}
	graphicsMod := graphics.FnModule{
		FnUpdateViewableChunks: func(viewableChunks map[graphics.ChunkEvent]struct{}) {
			actualViewableChunks = viewableChunks
		},
	}
	worldMod := world.New(graphicsMod)
	worldMod.UpdateView(map[world.ChunkEvent]struct{}{
		{1, 2, 3}: {},
	})

	if !reflect.DeepEqual(expectedViewableChunks, actualViewableChunks) {
		t.Fatalf("expected graphics modules to receive %v viewable chunks, but got %v", expectedViewableChunks, actualViewableChunks)
	}
}

func TestWorldLoadedChunkCount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc          string
		loadChunks    []world.ChunkEvent
		unloadChunks  []world.ChunkEvent
		expectedCount int
	}{
		{
			desc:          "world starts with no loaded chunks",
			expectedCount: 0,
		},
		{
			desc: "world should load one chunk",
			loadChunks: []world.ChunkEvent{
				{1, 2, 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should load and unload the same chunk",
			loadChunks: []world.ChunkEvent{
				{1, 2, 3},
			},
			unloadChunks: []world.ChunkEvent{
				{1, 2, 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should not load the same chunk twice",
			loadChunks: []world.ChunkEvent{
				{1, 2, 3},
				{1, 2, 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should not unload the same chunk twice",
			loadChunks: []world.ChunkEvent{
				{1, 2, 3},
			},
			unloadChunks: []world.ChunkEvent{
				{1, 2, 3},
				{1, 2, 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should load two different chunks",
			loadChunks: []world.ChunkEvent{
				{1, 2, 3},
				{4, 5, 6},
			},
			expectedCount: 2,
		},
		{
			desc: "world cannot unload a chunk if it has none",
			unloadChunks: []world.ChunkEvent{
				{1, 2, 3},
			},
			expectedCount: 0,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			graphicsMod := graphics.FnModule{}
			worldMod := world.New(graphicsMod)
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
