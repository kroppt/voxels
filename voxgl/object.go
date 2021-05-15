package voxgl

import (
	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/kroppt/gfx"
)

// Object is a renderable set of vertices.
type Object struct {
	program   gfx.Program
	vao       gfx.VAO
	position  mgl.Vec3
	scale     mgl.Vec3
	rotation  mgl.Quat
	camerapos mgl.Vec3
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
		program:   prog,
		vao:       *vao,
		position:  mgl.Vec3{0.0, 0.0, 0.0},
		scale:     mgl.Vec3{1.0, 1.0, 1.0},
		rotation:  mgl.QuatIdent(),
		camerapos: mgl.Vec3{0.0, 0.0, -5.0},
	}, nil
}

// Render generates an image of the object with OpenGL.
func (o *Object) Render() {
	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(o.position[0], o.position[1], o.position[2]))
	model = model.Mul4(mgl.Scale3D(o.scale[0], o.scale[1], o.scale[2]))
	model = model.Mul4(o.rotation.Mat4())
	o.program.UploadUniformMat4("model", model)

	view := mgl.Ident4()
	view = view.Mul4(mgl.Translate3D(o.camerapos[0], o.camerapos[1], o.camerapos[2]))
	o.program.UploadUniformMat4("view", view)

	proj := mgl.Ident4()
	proj = proj.Mul4(mgl.Perspective(mgl.DegToRad(45.0), 16.0/9.0, 0.1, 100.0))
	o.program.UploadUniformMat4("projection", proj)

	o.program.Bind()
	o.vao.Draw()
	o.program.Unbind()
}

// Translate adds the given position to the object.
// X, y, and z are the OpenGL coordinates to add to each of their respective
// dimensions.
func (o *Object) Translate(x, y, z float32) {
	o.position = o.position.Add(mgl.Vec3{x, y, z})
}

// Scale scales up or down the object by the given amounts.
// X, y, and z are the fraction to multiply the given scale by.
func (o *Object) Scale(x, y, z float32) {
	o.position = mgl.Vec3{o.position[0] * x, o.position[1] * y, o.position[2] * z}
}

func (o *Object) CameraTranslate(x, y, z float32) {
	o.camerapos = o.camerapos.Add(mgl.Vec3{x, y, z})
}

// Rotate adds the given rotation to the object.
// X, y, and z are the angles to rotate about each of their respective axis.
func (o *Object) Rotate(x, y, z float32) {
	xrot := mgl.HomogRotate3DX(mgl.DegToRad(x))
	yrot := mgl.HomogRotate3DY(mgl.DegToRad(y))
	zrot := mgl.HomogRotate3DZ(mgl.DegToRad(z))
	rotquat := mgl.Mat4ToQuat(xrot.Mul4(yrot).Mul4(zrot))
	o.rotation = o.rotation.Mul(rotquat)
}

// Destroy frees external resources.
func (o *Object) Destroy() {
	o.program.Destroy()
	o.vao.Destroy()
}
