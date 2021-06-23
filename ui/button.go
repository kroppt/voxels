package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*Button)(nil)

// Button stores the internal data associated with a button UI element.
type Button struct {
	program *gfx.Program
	vao     *gfx.VAO
	gfx     Gfx
	parent  *Background
}

func NewButton(gfx Gfx, parentComponent *Background, screenWidth, screenHeight int32) *Button {
	layout := []int32{2, 4}

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	btn := &Button{
		vao:    vao,
		gfx:    gfx,
		parent: parentComponent,
	}

	btn.ReloadPosition(screenWidth, screenHeight)

	return btn
}

func (btn *Button) ReloadPosition(screenWidth, screenHeight int32) {
	border := btn.GetBorder()
	width := btn.GetWidth()
	posTL := f32Point{float32(btn.parent.GetBorder() + border), float32(screenHeight - btn.parent.GetBorder() - border)}                                  // top-left
	posBL := f32Point{float32(btn.parent.GetBorder() + border), float32(screenHeight - btn.parent.GetBorder() - btn.parent.GetHeight() + border)}         // bottom-left
	posTR := f32Point{float32(btn.parent.GetBorder() + border + width), float32(screenHeight - btn.parent.GetBorder() - border)}                          // top-right
	posBR := f32Point{float32(btn.parent.GetBorder() + border + width), float32(screenHeight - btn.parent.GetBorder() - btn.parent.GetHeight() + border)} // bottom-right

	red := float32(0.0)   // red
	green := float32(1.0) // green
	blue := float32(0.0)  // blue
	alpha := float32(1.0) // alpha

	vertices := []float32{
		posTL.x, posTL.y, red, green, blue, alpha,
		posBL.x, posBL.y, red, green, blue, alpha,
		posTR.x, posTR.y, red, green, blue, alpha,
		posBR.x, posBR.y, red, green, blue, alpha,
	}

	btn.gfx.VAOLoad(btn.vao, vertices, gl.STATIC_DRAW)
}

// GetVAO returns the vertex array object associated with the button.
func (btn *Button) GetVAO() *gfx.VAO {
	return btn.vao
}

func (btn *Button) GetProgram() *gfx.Program {
	return btn.program
}

func (btn *Button) SetProgram(program *gfx.Program) {
	btn.program = program
}

func (btn *Button) GetBorder() int32 {
	if btn == nil {
		return int32(0)
	}
	return int32(10)
}

func (btn *Button) GetWidth() int32 {
	if btn == nil {
		return int32(0)
	}
	return int32(200)
}

func (btn *Button) GetTop() int32 {
	if btn == nil {
		return int32(0)
	}
	if btn.parent == nil {
		return btn.GetBorder()
	}
	return btn.parent.GetTop() + btn.GetBorder()
}

func (btn *Button) GetLeft() int32 {
	if btn == nil {
		return int32(0)
	}
	if btn.parent == nil {
		return btn.GetBorder()
	}
	return btn.parent.GetLeft() + btn.GetBorder()
}
