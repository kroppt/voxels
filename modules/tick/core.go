package tick

import "github.com/kroppt/voxels/modules/camera"

type core struct {
	cameraMod   camera.Interface
	currentTick int
}

func (c *core) getTick() int {
	return c.currentTick
}

func (c *core) advanceTick() {
	c.currentTick++
	c.cameraMod.Tick()
}
