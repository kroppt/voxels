package chunk

import "github.com/kroppt/voxels/modules/graphics"

type graphicsMod interface {
	UpdateDirection(graphics.DirectionEvent)
	ShowVoxel(graphics.VoxelEvent)
}

type core struct {
	graphicsMod graphicsMod
}

func (c core) init(chunkSize uint) {
	for i := uint(0); i < chunkSize*chunkSize; i++ {
		c.graphicsMod.ShowVoxel(graphics.VoxelEvent{})
	}
}

func (c core) updatePosition(posEvent PositionEvent) {
}
