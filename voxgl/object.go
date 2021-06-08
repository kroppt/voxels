// +build !test

package voxgl

import (
	"fmt"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// Object is a renderable set of vertices.
type Object struct {
	program  gfx.Program
	vao      gfx.VAO
	position glm.Vec3
	scale    glm.Vec3
	rotation glm.Quat
	texture1 *gfx.Texture3D
}

// NewObject returns a newly created Object with the given vertices.
func NewObject(program gfx.Program, vertices []float32, layout []int32) (*Object, error) {
	vao := gfx.NewVAO(gl.POINTS, layout)

	vao.Load(vertices, gl.STATIC_DRAW)

	return &Object{
		program:  program,
		vao:      *vao,
		position: glm.Vec3{0.0, 0.0, 0.0},
		scale:    glm.Vec3{1.0, 1.0, 1.0},
		rotation: glm.QuatIdent(),
	}, nil
}

func (o *Object) Add3DTexture(dx, dy, dz int32) {
	data := make([]byte, dx*dy*dz)
	texture1, err := gfx.NewTexture3D(dx, dy, dz, data, gl.RED, 1, 1)
	if err != nil {
		panic(fmt.Sprint(err))
	}
	o.texture1 = &texture1
}

func (o *Object) SetSpotOnTexture1(i, j, k int32) error {
	return o.texture1.SetPixel(gfx.Point3D{
		X: i,
		Y: j,
		Z: k,
	}, []byte{1}, false)
}

func (o *Object) SetManySpotsOnTexture1(size, height int32, data []byte) error {
	return o.texture1.SetPixelArea(0, 0, 0, size, height, size, data, false)
}

func (o *Object) SetData(data []float32) {
	err := o.vao.Load(data, gl.STATIC_DRAW)
	if err != nil {
		panic("failed to set data")
	}
}

// Render generates an image of the object with OpenGL.
func (o *Object) Render() {
	// sw := util.Start()
	o.program.Bind()
	o.texture1.Bind()
	o.vao.Draw()
	o.texture1.Unbind()
	o.program.Unbind()
	// gl.Finish()
	// sw.StopRecordAverage("individual voxel render")
}

// Translate adds the given position to the object.
// X, y, and z are the OpenGL coordinates to add to each of their respective
// dimensions.
func (o *Object) Translate(x, y, z float32) {
	o.position = o.position.Add(&glm.Vec3{x, y, z})
}

// Scale scales up or down the object by the given amounts.
// X, y, and z are the fraction to multiply the given scale by.
func (o *Object) Scale(x, y, z float32) {
	o.scale = glm.Vec3{o.scale[0] * x, o.scale[1] * y, o.scale[2] * z}
}

// Rotate adds the given rotation to the object.
// X, y, and z are the angles to rotate about each of their respective axis.
func (o *Object) Rotate(x, y, z float32) {
	xrot := glm.HomogRotate3DX(glm.DegToRad(x))
	yrot := glm.HomogRotate3DY(glm.DegToRad(y))
	zrot := glm.HomogRotate3DZ(glm.DegToRad(z))
	xyrot := xrot.Mul4(&yrot)
	xyzrot := xyrot.Mul4(&zrot)
	rotquat := glm.Mat4ToQuat(&xyzrot)
	o.rotation = o.rotation.Mul(&rotquat)
}

// Destroy frees external resources.
func (o *Object) Destroy() {
	o.program.Destroy()
	o.vao.Destroy()
	o.texture1.Destroy()
}
