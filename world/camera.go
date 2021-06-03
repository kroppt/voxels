package world

import (
	"github.com/engoengine/glm"
)

// Camera contains position and rotation information for the camera. UpdateView
// should be called whenever the camera is translated to rotated.
type Camera struct {
	pos glm.Vec3
	rot glm.Quat
}

// NewCamera returns a new camera.
func NewCamera() *Camera {
	return &Camera{
		pos: [3]float32{},
		rot: glm.QuatIdent(),
	}
}

// GetPosition returns the position of the camera.
func (c *Camera) GetPosition() glm.Vec3 {
	return c.pos.Mul(-1.0)
}

// SetPosition sets the position of the camera to pos
func (c *Camera) SetPosition(pos *glm.Vec3) {
	c.pos = pos.Mul(-1)
}

// Translate adds the given translation to the position of the camera.
func (c *Camera) Translate(diff *glm.Vec3) {
	c.pos = c.pos.Sub(diff)
}

// quatLookAtV is a fixed version of GLM's QuatLookAtV that accounts for Y direction
func quatLookAtV(eye, center, up *glm.Vec3) glm.Quat {
	// glm bug fix
	if *eye == *center {
		return glm.QuatIdent()
	}
	// Copied from GLM and uncommented 2 lines below
	cme := center.Sub(eye)
	direction := cme.Normalized()

	min1 := glm.Vec3{0, 0, -1}
	rotDir := glm.QuatBetweenVectors(&min1, &direction)

	// Uncommented these 2 lines
	right := direction.Cross(up)
	upp := right.Cross(&direction)

	dup := glm.Vec3{0, 1, 0}
	upCur := rotDir.Rotate(&dup)
	rotUp := glm.QuatBetweenVectors(&upCur, &upp)

	rotTarget := rotUp.Mul(&rotDir)
	return rotTarget.Inverse()
}

// LookAt rotates the camera to look at a specified point
func (c *Camera) LookAt(center *glm.Vec3) {
	// do nothing if target is current position
	negatedPos := c.pos.Mul(-1)
	if negatedPos == *center {
		return
	}
	up := c.GetLookUp()
	quat := quatLookAtV(&negatedPos, center, &up)
	c.rot = quat
}

// GetRotationQuat returns the quaternion the camera is rotated with.
func (c *Camera) GetRotationQuat() glm.Quat {
	return c.rot
}

// Rotate rotates the camera by degrees about the given axis.
func (c *Camera) Rotate(axis *glm.Vec3, deg float32) {
	rad := glm.DegToRad(deg)
	quat := glm.QuatRotate(rad, axis)
	c.rot = c.rot.Mul(&quat)
}

// GetViewMat returns the 4x4 matrix associated with the view represented by the
// camera.
func (c *Camera) GetViewMat() glm.Mat4 {
	view := glm.Ident4()
	cur := c.rot.Mat4()
	view = view.Mul4(&cur)
	pos := glm.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())
	view = view.Mul4(&pos)
	return view
}

// GetProjMat returns the 4x4 matrix associated with the projection represented
// by the camera.
func (c *Camera) GetProjMat() glm.Mat4 {
	proj := glm.Ident4()
	persp := glm.Perspective(glm.DegToRad(45.0), 16.0/9.0, 0.1, 100.0)
	proj = proj.Mul4(&persp)
	return proj
}

// GetLookForward returns the forward-looking direction vector
func (c *Camera) GetLookForward() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{0.0, 0.0, -1.0})
}

// GetLookBack returns the backwards-looking direction vector
func (c *Camera) GetLookBack() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{0.0, 0.0, 1.0})
}

// GetLookRight returns the right-looking direction vector
func (c *Camera) GetLookRight() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{1.0, 0.0, 0.0})
}

// GetLookLeft returns the left-looking direction vector
func (c *Camera) GetLookLeft() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{-1.0, 0.0, 0.0})
}

// GetLookUp returns the up-looking direction vector
func (c *Camera) GetLookUp() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{0.0, 1.0, 0.0})
}

// GetLookDown returns the up-looking direction vector
func (c *Camera) GetLookDown() glm.Vec3 {
	inverse := c.rot.Inverse()
	return inverse.Rotate(&glm.Vec3{0.0, -1.0, 0.0})
}
