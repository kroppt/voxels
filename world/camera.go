package world

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

// Camera contains position and rotation information for the camera. GetViewMat
// should be used with rendering.
type Camera struct {
	pos mgl.Vec3
	rot mgl.Quat
}

// NewCamera returns a new camera.
func NewCamera() *Camera {
	return &Camera{
		pos: [3]float32{},
		rot: mgl.QuatIdent(),
	}
}

// GetPosition returns the position of the camera.
func (c *Camera) GetPosition() mgl.Vec3 {
	return c.pos.Mul(-1.0)
}

// Translate adds the given translation to the position of the camera.
func (c *Camera) Translate(diff mgl.Vec3) {
	c.pos = c.pos.Sub(diff)
}

// GetRotationQuat returns the quaternion the camera is rotated with.
func (c *Camera) GetRotationQuat() mgl.Quat {
	return c.rot
}

// Rotate rotates the camera by degrees about the given axis.
func (c *Camera) Rotate(axis mgl.Vec3, deg float32) {
	rad := mgl.DegToRad(deg)
	quat := mgl.QuatRotate(rad, axis)
	c.rot = c.rot.Mul(quat)
}

// GetViewMat returns the 4x4 matrix associated with the view represented by the
// camera.
func (c *Camera) GetViewMat() mgl.Mat4 {
	view := mgl.Ident4()
	view = view.Mul4(c.rot.Mat4())
	view = view.Mul4(mgl.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z()))
	return view
}

func (c *Camera) GetLookForward() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, 0.0, -1.0})
}

func (c *Camera) GetLookBack() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, 0.0, 1.0})
}

func (c *Camera) GetLookRight() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{1.0, 0.0, 0.0})
}

func (c *Camera) GetLookLeft() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{-1.0, 0.0, 0.0})
}
