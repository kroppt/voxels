package chunk

import (
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod graphics.Interface
	settingsMod settings.Interface
}

func (c core) init() {
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	for x := int32(-renderDistance); x <= renderDistance; x++ {
		for y := int32(-renderDistance); y <= renderDistance; y++ {
			for z := int32(-renderDistance); z <= renderDistance; z++ {
				c.graphicsMod.ShowChunk(graphics.ChunkEvent{PositionX: x, PositionY: y, PositionZ: z})
			}
		}
	}
}

func (c core) updatePosition(posEvent PositionEvent) {
}
