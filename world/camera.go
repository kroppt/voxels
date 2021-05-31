package world

import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/kroppt/gfx"
)

// Camera contains position and rotation information for the camera. UpdateView
// should be called whenever the camera is translated to rotated.
type Camera struct {
	pos mgl.Vec3
	rot mgl.Quat
	ubo *gfx.BufferObject
}

// NewCamera returns a new camera.
func NewCamera() (*Camera, error) {
	cam := Camera{
		pos: [3]float32{},
		rot: mgl.QuatIdent(),
		ubo: gfx.NewBufferObject(),
	}
	var mat mgl.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	cam.ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	cam.ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)

	if err := cam.UpdateView(); err != nil {
		return nil, err
	}
	if err := cam.UpdateProj(); err != nil {
		return nil, err
	}

	return &cam, nil
}

func (c *Camera) Destroy() {
	c.ubo.Destroy()
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

// UpdateView sends an updated view matrix to a uniform buffer object
func (c *Camera) UpdateView() error {
	view := c.getViewMat()
	// offset of 0 because view is the first member of the struct
	err := c.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	return err
}

// UpdateProj sends an updated projection matrix to a uniform buffer object
func (c *Camera) UpdateProj() error {
	proj := mgl.Ident4()
	proj = proj.Mul4(mgl.Perspective(mgl.DegToRad(45.0), 16.0/9.0, 0.1, 100.0))
	// offset of sizeof(mat4) because proj is the second member of the struct
	err := c.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(proj)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	return err
}

// getViewMat returns the 4x4 matrix associated with the view represented by the
// camera.
func (c *Camera) getViewMat() mgl.Mat4 {
	view := mgl.Ident4()
	view = view.Mul4(c.rot.Mat4())
	view = view.Mul4(mgl.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z()))
	return view
}

// GetLookForward returns the forward-looking direction vector
func (c *Camera) GetLookForward() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, 0.0, -1.0})
}

// GetLookBack returns the backwards-looking direction vector
func (c *Camera) GetLookBack() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, 0.0, 1.0})
}

// GetLookRight returns the right-looking direction vector
func (c *Camera) GetLookRight() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{1.0, 0.0, 0.0})
}

// GetLookLeft returns the left-looking direction vector
func (c *Camera) GetLookLeft() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{-1.0, 0.0, 0.0})
}

// GetLookUp returns the up-looking direction vector
func (c *Camera) GetLookUp() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, 1.0, 0.0})
}

// GetLookDown returns the up-looking direction vector
func (c *Camera) GetLookDown() mgl.Vec3 {
	return c.rot.Inverse().Rotate(mgl.Vec3{0.0, -1.0, 0.0})
}
