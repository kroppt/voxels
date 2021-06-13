package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// UI is a struct.
type UI struct {
	vao     *gfx.VAO
	program *gfx.Program
}

// Gfx is an interface of the functions being provided.
type Gfx interface {
	NewVAO(mode uint32, layout []int32) *gfx.VAO
	VAOLoad(vao *gfx.VAO, data []float32, usage uint32) error
	NewShader(source string, shaderType uint32) (gfx.Shader, error)
	NewProgram(shaders ...gfx.Shader) (gfx.Program, error)
}

type f32Point struct {
	x float32
	y float32
}

// New returns a new ui.UI.
func New(gfx Gfx) (*UI, error) {
	layout := []int32{2}
	posTL := f32Point{-1, 1} // top-left
	posBL := f32Point{-1, 0} // bottom-left
	posTR := f32Point{1, 1}  // top-right
	posBR := f32Point{1, 0}  // bottom-right
	vertices := []float32{
		posTL.x, posTL.y,
		posBL.x, posBL.y,
		posTR.x, posTR.y,
		posBR.x, posBR.y,
	}

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	gfx.VAOLoad(vao, vertices, gl.STATIC_DRAW)

	vshad, err := gfx.NewShader(vertColShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshad, err := gfx.NewShader(fragColShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(vshad, fshad)
	if err != nil {
		return nil, err
	}

	ui := &UI{
		vao:     vao,
		program: &prog,
	}
	return ui, nil
}

// Render renders the object.
func (ui *UI) Render() {
	ui.program.Bind()
	ui.vao.Draw()
	ui.program.Unbind()
}

// Destroy destroys all members of the struct.
func (ui *UI) Destroy() {
	if ui.vao != nil {
		ui.vao.Destroy()
	}
	if ui.program != nil {
		ui.program.Destroy()
	}
}

const vertColShader = `
	#version 420 core

	layout (location = 0) in vec2 pos;

	void main() {
		gl_Position = vec4(pos, 0.0f, 1.0f);
	}
`

const fragColShader = `
	#version 330

	out vec4 frag_color;

	void main() {
		frag_color = vec4(1.0f, 0.0f, 0.0f, 1.0f);
	}
`
