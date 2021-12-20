package chunk_test

import (
	"reflect"
	"testing"

	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
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
	t.Parallel()

	t.Run("return is non-nil", func(t *testing.T) {
		t.Parallel()
		graphicsMod := graphics.FnModule{
			FnUpdatePlayerDirection: func(graphics.DirectionEvent) {
			},
			FnShowChunk: func(chunkEvent graphics.ChunkEvent) {
			},
		}
		settingsMod := settings.FnRepository{}

		mod := chunk.New(graphicsMod, settingsMod, 1)

		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})

	t.Run("panic on nil settingsMod", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		graphicsMod := graphics.FnModule{
			FnUpdatePlayerDirection: func(graphics.DirectionEvent) {
			},
			FnShowChunk: func(chunkEvent graphics.ChunkEvent) {
			},
		}

		chunk.New(graphicsMod, nil, 1)
	})

	t.Run("when chunk module is created, show chunks in the render distance", func(t *testing.T) {
		t.Parallel()
		expected := map[graphics.ChunkEvent]struct{}{}
		for x := int32(-2); x <= 2; x++ {
			for y := int32(-2); y <= 2; y++ {
				for z := int32(-2); z <= 2; z++ {
					expected[graphics.ChunkEvent{
						PositionX: x,
						PositionY: y,
						PositionZ: z,
					}] = struct{}{}
				}
			}
		}
		actual := map[graphics.ChunkEvent]struct{}{}
		graphicsMod := graphics.FnModule{
			FnUpdatePlayerDirection: func(graphics.DirectionEvent) {
			},
			FnShowChunk: func(chunkEvent graphics.ChunkEvent) {
				actual[chunkEvent] = struct{}{}
			},
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
		}

		chunk.New(graphicsMod, settingsMod, 1)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})

	t.Run("player position is 0, 0, 0 by default", func(t *testing.T) {
		t.Parallel()
		expected := graphics.ChunkEvent{
			PositionX: 0,
			PositionY: 0,
			PositionZ: 0,
		}
		var actual graphics.ChunkEvent
		graphicsMod := graphics.FnModule{
			FnUpdatePlayerDirection: func(graphics.DirectionEvent) {
			},
			FnShowChunk: func(chunkEvent graphics.ChunkEvent) {
				actual = chunkEvent
			},
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 0
			},
		}

		chunk.New(graphicsMod, settingsMod, 1)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})
}

func TestModuleUpdatePlayerPosition(t *testing.T) {
	t.Run("when player position is moved, new chunks are shown", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expected := map[graphics.ChunkEvent]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expected[graphics.ChunkEvent{
					PositionX: 3,
					PositionY: y,
					PositionZ: z,
				}] = struct{}{}
			}
		}
		graphicsMod := &graphics.FnModule{}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
		}
		chunkMod := chunk.New(graphicsMod, settingsMod, chunkSize)
		chunkMod.UpdatePlayerPosition(chunk.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actual := map[graphics.ChunkEvent]struct{}{}
		graphicsMod.FnShowChunk = func(chunkEvent graphics.ChunkEvent) {
			actual[chunkEvent] = struct{}{}
		}

		chunkMod.UpdatePlayerPosition(chunk.PositionEvent{
			X: chunkSize + 5,
			Y: 0,
			Z: 0,
		})

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})
}
