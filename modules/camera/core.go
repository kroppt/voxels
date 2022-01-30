package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/player"
)

type core struct {
	playerMod     player.Interface
	pos           mgl.Vec3
	rot           mgl.Quat
	moveForwards  bool
	moveLeft      bool
	moveBackwards bool
	moveRight     bool
}

func (c *core) tick() {
	moved := false
	forward := c.rot.Rotate(mgl.Vec3{0.0, 0.0, -1.0})
	if c.moveForwards {
		moved = true
		c.pos = c.pos.Add(forward)
	} else if c.moveBackwards {
		moved = true
		c.pos = c.pos.Add(mgl.Vec3{0, 0, 1})
	}
	if c.moveRight {
		moved = true
		c.pos = c.pos.Add(mgl.Vec3{1, 0, 0})
	} else if c.moveLeft {
		moved = true
		c.pos = c.pos.Add(mgl.Vec3{-1, 0, 0})
	}
	if moved {
		c.playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: c.pos.X(),
			Y: c.pos.Y(),
			Z: c.pos.Z(),
		})
	}
}

func (c *core) handleMovementEvent(evt MovementEvent) {
	switch evt.Direction {
	case MoveForwards:
		c.moveForwards = evt.Pressed
		if c.moveForwards && c.moveBackwards {
			c.moveBackwards = false
		}
	case MoveRight:
		c.moveRight = evt.Pressed
		if c.moveRight && c.moveLeft {
			c.moveLeft = false
		}
	case MoveBackwards:
		c.moveBackwards = evt.Pressed
		if c.moveBackwards && c.moveForwards {
			c.moveForwards = false
		}
	case MoveLeft:
		c.moveLeft = evt.Pressed
		if c.moveLeft && c.moveRight {
			c.moveRight = false
		}
	}
}

func (c *core) handleLookEvent(evt LookEvent) {
	radX := evt.Right
	quatX := mgl.QuatRotate(radX, mgl.Vec3{0, -1, 0})

	c.rot = quatX.Mul(c.rot)

	earAxis := c.rot.Rotate(mgl.Vec3{-1, 0, 0})

	radY := evt.Down
	quatY := mgl.QuatRotate(radY, earAxis)

	c.rot = quatY.Mul(c.rot)

	c.playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: c.rot,
	})
}
