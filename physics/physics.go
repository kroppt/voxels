package physics

import (
	"time"

	"github.com/engoengine/glm"
)

// Pos is a location away from the origin location with units meters.
type Pos struct {
	glm.Vec3
}

// Vel is a rate of change of position with units meters per second.
type Vel struct {
	glm.Vec3
}

// AsPosition computes the position from the origin that the velocity affects
// over the given dt.
func (v Vel) AsPosition(dt time.Duration) Pos {
	secs := float32(dt) / float32(time.Second)
	return Pos{glm.Vec3{v.X() * secs, v.Y() * secs, v.Z() * secs}}
}

// Acc is the rate of change of velocity with units meters per second squared.
type Acc struct {
	glm.Vec3
}

// AsVelocity computes the velocity from the origin that the acceleration
// affects over the given dt.
func (a Acc) AsVelocity(dt time.Duration) Vel {
	secs := float32(dt) / float32(time.Second)
	return Vel{glm.Vec3{a.X() * secs, a.Y() * secs, a.Z() * secs}}
}
