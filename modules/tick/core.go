package tick

import (
	"github.com/kroppt/voxels/modules/camera"
)

type core struct {
	cameraMod    camera.Interface
	timeMod      Time
	tickRateNano int64
	lastTickNano int64
	currentTick  int
}

func (c *core) getTick() int {
	return c.currentTick
}

func (c *core) advanceTick() {
	c.lastTickNano = c.timeMod.Now().UnixNano()
	c.currentTick++
	c.cameraMod.Tick()
}

func (c *core) isNextTickReady() bool {
	now := c.timeMod.Now().UnixNano()
	then := c.lastTickNano
	return now >= then+c.tickRateNano
}
