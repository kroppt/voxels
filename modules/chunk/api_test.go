package chunk_test

import (
	"testing"

	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

// Test Plan:
// [X] Change name of graphics mod Update{Position,Direction} to UpdatePlayer{Position,Direction}
// [X] Make interfaces for each module in their packages
// [X] On startup, renders chunks around player
// [X]  - Change graphics mod ShowVoxel to ShowChunk
// [X]  - Change graphics mod VoxelEvent to ChunkEvent
// [ ] Updates player position, renders different chunks around player
// [ ]  - Test to check new chunks are shown
// [ ]  - Test that old chunks are hidden

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		graphicsMod := graphics.FnModule{
			FnUpdatePlayerDirection: func(graphics.DirectionEvent) {
			},
			FnShowChunk: func(chunkEvent graphics.ChunkEvent) {
			},
		}

		mod := chunk.New(graphicsMod)

		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})

}
