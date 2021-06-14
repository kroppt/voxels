package ui

import (
	"errors"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// Element is an interface that represents something that is rendered like a UI element.
type Element interface {
	GetVAO() *gfx.VAO
	GetColor() *[4]float32
}

// UI is a struct all of the Elements that need to be rendered along with the OpenGL Program.
type UI struct {
	elements []Element
	program  *gfx.Program
	gfx      Gfx
}

// Gfx is an interface of the functions being provided.
type Gfx interface {
	NewVAO(mode uint32, layout []int32) *gfx.VAO
	VAOLoad(vao *gfx.VAO, data []float32, usage uint32) error
	VAODraw(vao *gfx.VAO)
	NewShader(source string, shaderType uint32) (gfx.Shader, error)
	NewProgram(shaders ...gfx.Shader) (gfx.Program, error)
	ProgramBind(program *gfx.Program)
	ProgramUnbind(program *gfx.Program)
	ProgramUploadUniform(program *gfx.Program, uniformName string, data ...float32)
}

type f32Point struct {
	x float32
	y float32
}

// New returns a new ui.UI.
func New(gfx Gfx) (*UI, error) {

	vbgshad, err := gfx.NewShader(vertElementShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fbgshad, err := gfx.NewShader(fragElementShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(vbgshad, fbgshad)
	if err != nil {
		return nil, err
	}

	elements := make([]Element, 0)

	ui := &UI{
		elements: elements,
		program:  &prog,
		gfx:      gfx,
	}

	return ui, nil
}

// AddElement adds the given UI Element to the rendered UI.
func (ui *UI) AddElement(element Element) error {
	if ui.elements == nil {
		return errors.New("elements is nil")
	}
	ui.elements = append(ui.elements, element)
	return nil
}

// Render renders the object.
func (ui *UI) Render() {
	ui.gfx.ProgramBind(ui.program)
	for _, element := range ui.elements {
		ui.gfx.VAODraw(element.GetVAO())
		ui.gfx.ProgramUploadUniform(ui.program, "color", (*element.GetColor())[:]...)
	}
	ui.gfx.ProgramUnbind(ui.program)
}

// Destroy destroys all members of the struct.
func (ui *UI) Destroy() {
	for _, element := range ui.elements {
		element.GetVAO().Destroy()
	}
	if ui.program != nil {
		ui.program.Destroy()
	}
}

const vertElementShader = `
	#version 420 core

	layout (location = 0) in vec2 pos;

	void main() {
		gl_Position = vec4(pos, 0.0f, 1.0f);
	}
`

const fragElementShader = `
	#version 330

	out vec4 frag_color;

	uniform vec4 color;

	void main() {
		frag_color = vec4(color[0], color[1], color[2], color[3]);
	}
`
