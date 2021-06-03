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
}

// NewObject returns a newly created Object with the given vertices.
func NewObject(vertShader string, fragShader string, vertices [][]float32, layout []int32) (*Object, error) {
	vshad, err := gfx.NewShader(vertShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fshad, err := gfx.NewShader(fragShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	prog, err := gfx.NewProgram(vshad, fshad)
	if err != nil {
		return nil, err
	}

	vao := gfx.NewVAO(gl.TRIANGLES, layout)

	points := []float32{}
	for _, v := range vertices {
		points = append(points, v[:]...)
	}

	err = vao.Load(points, gl.STATIC_DRAW)
	if err != nil {
		return nil, err
	}

	return &Object{
		program:  prog,
		vao:      *vao,
		position: glm.Vec3{0.0, 0.0, 0.0},
		scale:    glm.Vec3{1.0, 1.0, 1.0},
		rotation: glm.QuatIdent(),
	}, nil
}

// Render generates an image of the object with OpenGL.
func (o *Object) Render() {
	model := glm.Ident4()
	trans := glm.Translate3D(o.position[0], o.position[1], o.position[2])
	model = model.Mul4(&trans)
	scale := glm.Scale3D(o.scale[0], o.scale[1], o.scale[2])
	model = model.Mul4(&scale)
	rot := o.rotation.Mat4()
	model = model.Mul4(&rot)
	err := o.program.UploadUniformMat4("model", model)
	if err != nil {
		panic(fmt.Errorf("error uploading uniform \"model\": %w", err))
	}

	o.program.Bind()
	o.vao.Draw()
	o.program.Unbind()
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
}
