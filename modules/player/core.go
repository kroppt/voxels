package player

import "github.com/kroppt/voxels/modules/chunk"

type chunkMod interface {
	UpdatePosition(chunk.PositionEvent)
}

type core struct {
	chunkMod chunkMod
	x        int32
	y        int32
	z        int32
}

func (c *core) handleMovementEvent(evt MovementEvent) {
	switch evt.Direction {
	case MoveForwards:
		c.z--
	case MoveRight:
		c.x++
	case MoveBackwards:
		c.z++
	case MoveLeft:
		c.x--
	}
	posEvent := chunk.PositionEvent{
		X: c.x,
		Y: c.y,
		Z: c.z,
	}
	c.chunkMod.UpdatePosition(posEvent)
}

func (c *core) handleLookEvent(evt LookEvent) {
}
