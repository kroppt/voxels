package chunk_test

import (
	"fmt"
	"testing"

	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type fnGraphicsMod struct {
	fnUpdateDirection func(graphics.DirectionEvent)
	fnShowVoxel       func(graphics.VoxelEvent)
}

func (fn fnGraphicsMod) UpdateDirection(directionEvent graphics.DirectionEvent) {
	fn.fnUpdateDirection(directionEvent)
}

func (fn fnGraphicsMod) ShowVoxel(voxelEvent graphics.VoxelEvent) {
	fn.fnShowVoxel(voxelEvent)
}

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		graphicsMod := fnGraphicsMod{
			fnUpdateDirection: func(graphics.DirectionEvent) {
			},
			fnShowVoxel: func(voxelEvent graphics.VoxelEvent) {
			},
		}

		mod := chunk.New(graphicsMod, 1)

		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})

	t.Run("calls ShowVoxel based on chunk size", func(t *testing.T) {
		testCases := []struct {
			chunkSize uint
			expect    uint
		}{
			{
				chunkSize: 1,
				expect:    1 * 1,
			},
			{
				chunkSize: 2,
				expect:    2 * 2,
			},
			{
				chunkSize: 3,
				expect:    3 * 3,
			},
		}
		for _, tC := range testCases {
			var calls uint
			t.Run(fmt.Sprintf("called %v times for chunk size %v", calls, tC.chunkSize), func(t *testing.T) {
				graphicsMod := fnGraphicsMod{
					fnUpdateDirection: func(graphics.DirectionEvent) {
					},
					fnShowVoxel: func(graphics.VoxelEvent) {
						calls++
					},
				}

				_ = chunk.New(graphicsMod, tC.chunkSize)

				if calls != tC.expect {
					t.Fatalf("expected %v calls but got %v", tC.expect, calls)
				}
			})
		}
	})

	t.Run("zero chunk size causes panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic but got none")
			}
		}()

		graphicsMod := fnGraphicsMod{
			fnUpdateDirection: func(graphics.DirectionEvent) {
			},
			fnShowVoxel: func(graphics.VoxelEvent) {
			},
		}

		_ = chunk.New(graphicsMod, 0)
	})
}
