package graphics

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type CrosshairObject struct {
	prog gfx.Program
	vao  *gfx.VAO
}

func NewCrosshairObject(size, aspect float32) (*CrosshairObject, error) {
	vao := gfx.NewVAO(gl.LINES, []int32{2, 4})
	vertices := []float32{
		-size / aspect, 0, 0.0, 1.0, 1.0, 1.0,
		size / aspect, 0, 1.0, 1.0, 0.0, 1.0,
		0, -size, 1.0, 0.0, 1.0, 1.0,
		0, size, 0.0, 1.0, 0.0, 1.0,
	}
	vao.Load(vertices, gl.STATIC_DRAW)

	vshad, err := gfx.NewShader(vertCrossShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshad, err := gfx.NewShader(fragCrossShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(vshad, fshad)
	if err != nil {
		return nil, err
	}

	prog.UploadUniform("aspect", aspect)

	return &CrosshairObject{
		prog: prog,
		vao:  vao,
	}, nil
}

func (c *CrosshairObject) Render() {
	c.prog.Bind()
	c.vao.Draw()
	c.prog.Unbind()
}

func (c *CrosshairObject) Destroy() {
	c.vao.Destroy()
	c.prog.Destroy()
}

const vertCrossShader = `
	#version 420 core

	layout (location = 0) in vec2 pos;
	layout (location = 1) in vec4 col;

	out Vertex {
		vec4 col;
	} OUT;

	void main()
	{
		gl_Position = vec4(pos.xy, 0.0, 1.0);
		OUT.col = col;
	}
`

const fragCrossShader = `
	#version 400

	in Vertex {
		vec4 col;
	} IN;

	out vec4 frag_color;

	void main() {
		frag_color = IN.col;
	}
`
