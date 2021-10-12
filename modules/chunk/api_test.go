package chunk_test

import (
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

		mod := chunk.New(graphicsMod)

		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})

	t.Run("calls ShowVoxel", func(t *testing.T) {
		var evt *graphics.VoxelEvent
		graphicsMod := fnGraphicsMod{
			fnUpdateDirection: func(graphics.DirectionEvent) {
			},
			fnShowVoxel: func(voxelEvent graphics.VoxelEvent) {
				evt = &voxelEvent
			},
		}

		_ = chunk.New(graphicsMod)

		if evt == nil {
			t.Fatal("expected function to be called")
		}
	})
}
