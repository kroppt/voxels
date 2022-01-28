package camera

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
)

type core struct {
	playerMod   player.Interface
	graphicsMod graphics.Interface
	x           int32
	y           int32
	z           int32
	rot         mgl.Quat
}

func (c *core) handleMovementEvent(evt MovementEvent) {
	if evt.Pressed {
		return
	}
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
	posEvent := player.PositionEvent{
		X: c.x,
		Y: c.y,
		Z: c.z,
	}
	c.playerMod.UpdatePlayerPosition(posEvent)
}

func (c *core) handleLookEvent(evt LookEvent) {
	radX := evt.Right
	quatX := mgl.QuatRotate(radX, mgl.Vec3{0, -1, 0})

	c.rot = quatX.Mul(c.rot)

	downAxis := c.rot.Rotate(mgl.Vec3{-1, 0, 0})

	radY := evt.Down
	quatY := mgl.QuatRotate(radY, downAxis)

	c.rot = quatY.Mul(c.rot)

	c.playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: c.rot,
	})
}