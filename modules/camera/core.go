package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/player"
)

type core struct {
	playerMod player.Interface
	pos       mgl.Vec3
	rot       mgl.Quat
	// 0=foward,1=backward,2=left,3=right,4=down,5=up
	keyPressed   [6]bool
	moveNextTick [6]bool
}

func (c *core) tick() {
	moved := false
	if c.moveNextTick[0] {
		if !c.keyPressed[0] {
			c.moveNextTick[0] = false
		}
		moved = true
		c.pos = c.pos.Add(c.rot.Rotate(mgl.Vec3{0.0, 0.0, -1.0}))
	} else if c.moveNextTick[1] {
		if !c.keyPressed[1] {
			c.moveNextTick[1] = false
		}
		moved = true
		c.pos = c.pos.Add(c.rot.Rotate(mgl.Vec3{0.0, 0.0, 1.0}))
	}
	if c.moveNextTick[3] {
		if !c.keyPressed[3] {
			c.moveNextTick[3] = false
		}
		moved = true
		c.pos = c.pos.Add(c.rot.Rotate(mgl.Vec3{1.0, 0.0, 0.0}))
	} else if c.moveNextTick[2] {
		if !c.keyPressed[2] {
			c.moveNextTick[2] = false
		}
		moved = true
		c.pos = c.pos.Add(c.rot.Rotate(mgl.Vec3{-1.0, 0.0, 0.0}))
	}
	if c.moveNextTick[5] {
		if !c.keyPressed[5] {
			c.moveNextTick[5] = false
		}
		moved = true
		c.pos = c.pos.Add(mgl.Vec3{0.0, 1.0, 0.0})
	} else if c.moveNextTick[4] {
		if !c.keyPressed[4] {
			c.moveNextTick[4] = false
		}
		moved = true
		c.pos = c.pos.Add(mgl.Vec3{0.0, -1.0, 0.0})
	}
	if moved {
		c.playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: c.pos.X(),
			Y: c.pos.Y(),
			Z: c.pos.Z(),
		})
	}
}

func (c *core) handleKeyPressFlags(idx int, pressed bool) {
	diff := 1
	if idx%2 == 1 {
		diff = -1
	}
	c.keyPressed[idx] = pressed
	if c.keyPressed[idx] {
		c.moveNextTick[idx] = true
		c.moveNextTick[idx+diff] = false
	}
	if c.keyPressed[idx] && c.keyPressed[idx+diff] {
		c.keyPressed[idx+diff] = false
	}
}

func (c *core) handleMovementEvent(evt MovementEvent) {
	switch evt.Direction {
	case MoveForwards:
		c.handleKeyPressFlags(0, evt.Pressed)
	case MoveBackwards:
		c.handleKeyPressFlags(1, evt.Pressed)
	case MoveLeft:
		c.handleKeyPressFlags(2, evt.Pressed)
	case MoveRight:
		c.handleKeyPressFlags(3, evt.Pressed)
	case MoveDown:
		c.handleKeyPressFlags(4, evt.Pressed)
	case MoveUp:
		c.handleKeyPressFlags(5, evt.Pressed)
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
