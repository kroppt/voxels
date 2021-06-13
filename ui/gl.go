package ui

import "github.com/kroppt/gfx"

type GlGfx struct {
}

func (g *GlGfx) NewVAO(mode uint32, layout []int32) *gfx.VAO {
	return gfx.NewVAO(mode, layout)
}

func (g *GlGfx) VAOLoad(vao *gfx.VAO, data []float32, usage uint32) error {
	return vao.Load(data, usage)
}

func (g *GlGfx) NewShader(source string, shaderType uint32) (gfx.Shader, error) {
	return gfx.NewShader(source, shaderType)
}

func (g *GlGfx) NewProgram(shaders ...gfx.Shader) (gfx.Program, error) {
	return gfx.NewProgram(shaders...)
}
