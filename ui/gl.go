package ui

import (
	"github.com/kroppt/gfx"
)

var _ Gfx = (*GlGfx)(nil)

type GlGfx struct {
}

func (g *GlGfx) NewVAO(mode uint32, layout []int32) *gfx.VAO {
	return gfx.NewVAO(mode, layout)
}

func (g *GlGfx) VAOLoad(vao *gfx.VAO, data []float32, usage uint32) error {
	return vao.Load(data, usage)
}

func (g *GlGfx) VAODraw(vao *gfx.VAO) {
	vao.Draw()
}

func (g *GlGfx) NewShader(source string, shaderType uint32) (gfx.Shader, error) {
	return gfx.NewShader(source, shaderType)
}

func (g *GlGfx) NewProgram(shaders ...gfx.Shader) (gfx.Program, error) {
	return gfx.NewProgram(shaders...)
}

func (g *GlGfx) ProgramUploadUniform(program *gfx.Program, uniformName string, data ...float32) {
	program.UploadUniform(uniformName, data...)
}

func (g *GlGfx) ProgramBind(program *gfx.Program) {
	program.Bind()
}

func (g *GlGfx) ProgramUnbind(program *gfx.Program) {
	program.Unbind()
}
