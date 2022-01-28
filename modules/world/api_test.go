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
