package ui

import (
	"errors"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

const fontSize = 16

// Element is an interface that represents something that is rendered like a UI element.
type Element interface {
	GetProgram() *gfx.Program
	SetProgram(program *gfx.Program)
	GetVAO() *gfx.VAO
	ReloadPosition(screenWidth, screenHeight int32)
}

// UI is a struct all of the Elements that need to be rendered along with the OpenGL Program.
type UI struct {
	elements    []Element
	program     *gfx.Program
	textProgram *gfx.Program
	gfx         Gfx
	escapeMenu  bool
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
	ProgramUploadUniform(program *gfx.Program, uniformName string, data ...float32) error
	LoadFontTexture(fontName string, fontSize int32) (*gfx.FontInfo, error)
}

type f32Point struct {
	x float32
	y float32
}

// New returns a new ui.UI.
func New(gfx Gfx) (*UI, error) {

	vshad, err := gfx.NewShader(vertElementShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshad, err := gfx.NewShader(fragElementShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(vshad, fshad)
	if err != nil {
		return nil, err
	}

	vtextshad, err := gfx.NewShader(glyphShaderVertex, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	ftextshad, err := gfx.NewShader(glyphShaderFragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	textProg, err := gfx.NewProgram(vtextshad, ftextshad)

	elements := make([]Element, 0)

	ui := &UI{
		elements:    elements,
		program:     &prog,
		textProgram: &textProg,
		gfx:         gfx,
	}

	ui.gfx.ProgramUploadUniform(ui.program, "screenSize", float32(1920), float32(1080))

	fnt, err := gfx.LoadFontTexture("NotoMono-Regular.ttf", fontSize)
	if err != nil {
		return nil, err
	}

	err = ui.gfx.ProgramUploadUniform(ui.textProgram, "screen_size", float32(1920), float32(1080))
	if err != nil {
		return nil, err
	}

	textColor := [4]float32{0.0, 0.0, 0.0, 1.0}

	err = ui.gfx.ProgramUploadUniform(ui.textProgram, "tex_size", float32(fnt.GetTexture().GetWidth()),
		float32(fnt.GetTexture().GetHeight()))
	if err != nil {
		return nil, err
	}
	err = ui.gfx.ProgramUploadUniform(ui.textProgram, "text_color", textColor[0], textColor[1], textColor[2], textColor[3])
	if err != nil {
		return nil, err
	}

	return ui, nil
}

// AddElement adds the given UI Element to the rendered UI.
func (ui *UI) AddElement(element Element) error {
	if ui.elements == nil {
		return errors.New("elements is nil")
	}
	_, ok := element.(*Text)
	if ok {
		element.SetProgram(ui.textProgram)
	} else {
		element.SetProgram(ui.program)
	}
	ui.elements = append(ui.elements, element)
	return nil
}

// Render renders the object.
func (ui *UI) Render() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	for _, element := range ui.elements {
		_, okEsc := element.(*EscapeMenu)
		if okEsc && !ui.escapeMenu {
			// skip rendering the EscapeMenu
			continue
		}

		ui.gfx.ProgramBind(element.GetProgram())
		text, okText := element.(*Text)
		if okText {
			text.FontTextureBind()
		}
		ui.gfx.VAODraw(element.GetVAO())
		if okText {
			text.FontTextureUnbind()
		}
		ui.gfx.ProgramUnbind(element.GetProgram())
	}
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
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

// ToggleEscapeMenu enables, or disables if it is already enabled, the escape menu.
func (ui *UI) ToggleEscapeMenu() {
	if ui != nil {
		ui.escapeMenu = !ui.escapeMenu
	}
}

const vertElementShader = `
	#version 420 core

	layout (location = 0) in vec2 pixelPos;
	layout (location = 1) in vec4 color;

	out vec4 vertexColor;

	// (width, height)
	uniform vec2 screenSize;

	void main() {
		vec2 pos = 2.0f * vec2(pixelPos[0] / screenSize[0], pixelPos[1] / screenSize[1]) - 1.0f;
		gl_Position = vec4(pos, 0.0f, 1.0f);
		vertexColor = color;
	}
`

const fragElementShader = `
	#version 330

	in vec4 vertexColor;

	out vec4 frag_color;

	void main() {
		frag_color = vec4(vertexColor);
	}
`

// Uniform `tex_size` is the (width, height) of the texture.
// Input `position_in` is typical openGL position coordinates.
// Input `tex_pixels` is the (x, y) of the vertex in the texture starting
// at (left, top).
// Output `tex_coords` is typical texture coordinates for fragment shader.
const glyphShaderVertex = `
	#version 330
	uniform vec2 tex_size;
	uniform vec2 screen_size;
	layout(location = 0) in vec2 position_in;
	layout(location = 1) in vec2 tex_pixels;
	out vec2 tex_coords;
	void main() {
		vec2 glSpace = vec2(2.0, 2.0) * (position_in / screen_size) + vec2(-1.0, -1.0);
		gl_Position = vec4(glSpace, 0.0, 1.0);
		tex_coords = vec2(tex_pixels.x / tex_size.x, tex_pixels.y / tex_size.y);
	}
`

const glyphShaderFragment = `
	#version 330
	uniform sampler2D frag_tex;
	uniform vec4 text_color;
	in vec2 tex_coords;
	out vec4 frag_color;
	void main() {
		frag_color = vec4(text_color.xyz, texture(frag_tex, tex_coords).r * text_color.w);
	}
`
