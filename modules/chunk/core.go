package chunk

import "github.com/kroppt/voxels/modules/graphics"

type graphicsMod interface {
	UpdateDirection(graphics.DirectionEvent)
	ShowVoxel(graphics.VoxelEvent)
}

type core struct {
	graphicsMod graphicsMod
}

func (c core) init() {
	c.graphicsMod.ShowVoxel(graphics.VoxelEvent{})
}

func (c core) updatePosition(posEvent PositionEvent) {
}
