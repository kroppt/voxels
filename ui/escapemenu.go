package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*EscapeMenu)(nil)

// EscapeMenu stores the information needed for an escape menu Element.
type EscapeMenu struct {
	program *gfx.Program
	vao     *gfx.VAO
	gfx     Gfx
}

// NewEscapeMenu creates and returns a pointer to a new EscapeMenu.
func NewEscapeMenu(gfx Gfx, screenWidth, screenHeight int32) *EscapeMenu {
	layout := []int32{2, 4}

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	menu := &EscapeMenu{
		vao: vao,
		gfx: gfx,
	}

	menu.ReloadPosition(screenWidth, screenHeight)

	return menu
}

// GetProgram returns the stored program.
func (menu *EscapeMenu) GetProgram() *gfx.Program {
	return menu.program
}

// SetProgram sets the stored program.
func (menu *EscapeMenu) SetProgram(program *gfx.Program) {
	menu.program = program
}

// GetVAO returns the stored VAO.
func (menu *EscapeMenu) GetVAO() *gfx.VAO {
	return menu.vao
}

// GetWidth returns the width of the element in pixels.
func (menu *EscapeMenu) GetWidth() int32 {
	if menu == nil {
		return int32(0)
	}
	return int32(300)
}

// GetHeight returns the height of the element in pixels.
func (menu *EscapeMenu) GetHeight() int32 {
	if menu == nil {
		return int32(0)
	}
	return int32(400)
}

// ReloadPosition loads the VAO with the position to make this element centered.
func (menu *EscapeMenu) ReloadPosition(screenWidth, screenHeight int32) {
	center := [2]int32{screenWidth / 2, screenHeight / 2}

	width := menu.GetWidth()
	height := menu.GetHeight()
	posTL := f32Point{float32(center[0] - (width / 2)), float32(center[1] + (height / 2))} // top-left
	posBL := f32Point{float32(center[0] - (width / 2)), float32(center[1] - (height / 2))} // bottom-left
	posTR := f32Point{float32(center[0] + (width / 2)), float32(center[1] + (height / 2))} // top-right
	posBR := f32Point{float32(center[0] + (width / 2)), float32(center[1] - (height / 2))} // bottom-right

	var red float32 = 1.0   // red
	var green float32 = 0.0 // green
	var blue float32 = 0.0  // blue
	var alpha float32 = 1.0 // alpha

	vertices := []float32{
		posTL.x, posTL.y, red, green, blue, alpha,
		posBL.x, posBL.y, red, green, blue, alpha,
		posTR.x, posTR.y, red, green, blue, alpha,
		posBR.x, posBR.y, red, green, blue, alpha,
	}

	menu.gfx.VAOLoad(menu.vao, vertices, gl.STATIC_DRAW)
}
