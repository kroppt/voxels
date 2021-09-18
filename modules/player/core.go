package player

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type core struct {
	chunkMod    chunkMod
	graphicsMod graphicsMod
	x           int32
	y           int32
	z           int32
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
	rotX := mgl.QuatIdent()
	radX := evt.Right
	quatX := mgl.QuatRotate(radX, mgl.Vec3{0, -1, 0})
	rotX = rotX.Mul(quatX)

	rotY := mgl.QuatIdent()
	radY := evt.Down
	quatY := mgl.QuatRotate(radY, mgl.Vec3{-1, 0, 0})
	rotY = rotY.Mul(quatY)

	c.graphicsMod.UpdateDirection(graphics.DirectionEvent{
		Rotation: rotX.Mul(rotY),
	})
}
