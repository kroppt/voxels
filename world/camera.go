package world

import (
	"math"

	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
)

// Camera contains position and rotation information for the camera. UpdateView
// should be called whenever the camera is translated to rotated.
type Camera struct {
	pos     glm.Vec3
	rot     glm.Quat
	near    float32
	far     float32
	aspect  float32
	fovyDeg float32
	dirty   bool
}

// NewCameraDefault returns a new camera with default settings.
func NewCameraDefault() *Camera {
	return &Camera{
		fovyDeg: 60.0,
		aspect:  16.0 / 9.0,
		near:    0.1,
		far:     100.0,
		pos:     [3]float32{},
		rot:     glm.QuatIdent(),
	}
}

// NewCameraCustom returns a new camera with custom projection settings
func NewCameraCustom(fovyDeg, aspect, near, far float32) *Camera {
	return &Camera{
		fovyDeg: fovyDeg,
		aspect:  aspect,
		near:    near,
		far:     far,
		pos:     [3]float32{},
		rot:     glm.QuatIdent(),
	}
}

func (c *Camera) IsDirty() bool {
	return c.dirty
}

func (c *Camera) Clean() {
	c.dirty = false
}

func (c *Camera) GetNear() float32 {
	return c.near
}
func (c *Camera) GetFar() float32 {
	return c.far
}
func (c *Camera) GetFovy() float32 {
	return c.fovyDeg
}
func (c *Camera) GetAspect() float32 {
	return c.aspect
}

// GetPosition returns the position of the camera.
func (c *Camera) GetPosition() glm.Vec3 {
	return c.pos.Mul(-1.0)
}

// SetPosition sets the position of the camera to pos
func (c *Camera) SetPosition(pos *glm.Vec3) {
	c.pos = pos.Mul(-1)
	c.dirty = true
}

// Translate adds the given translation to the position of the camera.
func (c *Camera) Translate(diff *glm.Vec3) {
	c.pos = c.pos.Sub(diff)
	c.dirty = true
}

func (c *Camera) AsVoxelPos() VoxelPos {
	// negated
	pos := c.GetPosition()
	if pos.X() < 0 {
		pos[0] -= 1
	}
	if pos.Y() < 0 {
		pos[1] -= 1
	}
	if pos.Z() < 0 {
		pos[2] -= 1
	}
	return VoxelPos{
		int(pos.X()),
		int(pos.Y()),
		int(pos.Z()),
	}
}

type FRange struct {
	Min   float32
	Max   float32
	delta float32
}

type WorldRange struct {
	X FRange
	Y FRange
	Z FRange
}

func (rng WorldRange) ForEach(fn func(glm.Vec3)) {
	for x := rng.X.Min; x <= rng.X.Max; x += rng.X.delta {
		for y := rng.Y.Min; y <= rng.Y.Max; y += rng.Y.delta {
			for z := rng.Z.Min; z <= rng.Z.Max; z += rng.Z.delta {
				fn(glm.Vec3{x, y, z})
			}
		}
	}
}

// IsWithinFrustum returns whether a specified AABB is within the frustum
// view of the camera.
func (c *Camera) IsWithinFrustum(corner glm.Vec3, dx, dy, dz float32) bool {
	eye := c.GetPosition()
	dir := c.GetLookForward()
	left := c.GetLookLeft()
	right := c.GetLookRight()
	up := c.GetLookUp()
	down := c.GetLookDown()
	// far plane math
	farDist := dir.Mul(c.far)
	farCenter := eye.Add(&farDist)
	fovyRad := glm.DegToRad(c.fovyDeg / 2.0)
	fhh := c.far * float32(math.Tan(float64(fovyRad)))
	fhw := c.aspect * fhh
	farLeftOff := left.Mul(fhw)
	farRightOff := right.Mul(fhw)
	farUpOff := up.Mul(fhh)
	farDownOff := down.Mul(fhh)
	ftl := farCenter.Add(&farLeftOff)
	ftl = ftl.Add(&farUpOff)
	fbl := farCenter.Add(&farLeftOff)
	fbl = fbl.Add(&farDownOff)
	ftr := farCenter.Add(&farRightOff)
	ftr = ftr.Add(&farUpOff)
	fbr := farCenter.Add(&farRightOff)
	fbr = fbr.Add(&farDownOff)
	// near plane math
	nearDist := dir.Mul(c.near)
	nearCenter := eye.Add(&nearDist)
	nhh := c.near * float32(math.Tan(float64(fovyRad/2.0)))
	nhw := c.aspect * nhh
	nearLeftOff := left.Mul(nhw)
	nearUpOff := up.Mul(nhh)
	nleft := nearCenter.Add(&nearLeftOff)
	nup := nearCenter.Add(&nearUpOff)

	planeTriangles := [6][3]glm.Vec3{
		{eye, ftl, fbl},          // left
		{eye, ftr, ftl},          // top
		{eye, fbr, ftr},          // right
		{eye, fbl, fbr},          // bottom
		{fbl, ftl, ftr},          // far
		{nearCenter, nup, nleft}, // near
	}
	cubeRange := WorldRange{
		X: FRange{corner.X(), corner.X() + dx, dx},
		Y: FRange{corner.Y(), corner.Y() + dy, dy},
		Z: FRange{corner.Z(), corner.Z() + dz, dz},
	}
	for _, tri := range planeTriangles {
		in := 0
		cubeRange.ForEach(func(v glm.Vec3) {
			// every corner of cube
			if !geo.PointOutsidePlane(&v, &tri[0], &tri[1], &tri[2]) {
				in++
				// TODO break early, change ForEach
			}
		})
		if in == 0 {
			return false
		}
	}
	return true
}

// quatLookAtV is a fixed version of GLM's QuatLookAtV that accounts for Y direction
func quatLookAtV(eye, center, up, forward *glm.Vec3) glm.Quat {
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
	// the direction of target is parallel to up: either opposite or same direction
	if upp == (glm.Vec3{0, 0, 0}) {
		if direction == *up {
			upp = forward.Mul(-1)
		} else if direction == up.Mul(-1) {
			upp = *forward
		} else {
			panic("unexpected rounding error")
		}
	}

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
	forward := c.GetLookForward()
	quat := quatLookAtV(&negatedPos, center, &up, &forward)
	c.rot = quat
	c.dirty = true
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
	c.dirty = true
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
	persp := glm.Perspective(glm.DegToRad(c.fovyDeg), c.aspect, c.near, c.far)
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
