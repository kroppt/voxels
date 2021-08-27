package player

import (
	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
)

type chunkMod interface {
	UpdatePosition(chunk.PositionEvent)
}

type graphicsMod interface {
	UpdateDirection(graphics.DirectionEvent)
}

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
	rotX := glm.QuatIdent()
	radX := evt.Right
	quatX := glm.QuatRotate(radX, &glm.Vec3{0, -1, 0})
	rotX = rotX.Mul(&quatX)

	rotY := glm.QuatIdent()
	radY := evt.Down
	right := rotX.Rotate(&glm.Vec3{-1, 0, 0})
	quatY := glm.QuatRotate(radY, &right)
	rotY = rotY.Mul(&quatY)

	/*
		rotY := glm.QuatIdent()
		radY := evt.Down
		quatY := glm.QuatRotate(radY, &glm.Vec3{-1, 0, 0})
		rotY = rotY.Mul(&quatY)

		rotX := glm.QuatIdent()
		radX := evt.Right
		down := rotY.Rotate(&glm.Vec3{0, -1, 0})
		quatX := glm.QuatRotate(radX, &down)
		rotX = rotX.Mul(&quatX)
	*/

	/*
		rotX := glm.QuatIdent()
		radX := evt.Right
		quatX := glm.QuatRotate(radX, &glm.Vec3{0, -1, 0})
		rotX = rotX.Mul(&quatX)

		rotY := glm.QuatIdent()
		radY := evt.Down
		quatY := glm.QuatRotate(radY, &glm.Vec3{-1, 0, 0})
		rotY = rotY.Mul(&quatY)
	*/

	c.graphicsMod.UpdateDirection(graphics.DirectionEvent{
		Rotation: rotX.Mul(&rotY),
	})
}
