package settings

type core struct {
	fovY   float32
	width  int32
	height int32
}

func (c *core) setFOV(degY float32) {
	c.fovY = degY
}

func (c *core) getFOV() float32 {
	return c.fovY
}

func (c *core) setResolution(width, height int32) {
	c.width = width
	c.height = height
}

func (c *core) getResolution() (int32, int32) {
	return c.width, c.height
}
