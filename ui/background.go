package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*Background)(nil)

// Background stores the internal data associated with a background UI element.
type Background struct {
	vao *gfx.VAO
	gfx Gfx
}

// NewBackground creates a Background element.
func NewBackground(gfx Gfx, screenWidth, screenHeight int32) *Background {
	layout := []int32{2, 4}

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	bg := &Background{
		vao: vao,
		gfx: gfx,
	}

	bg.ReloadPosition(screenWidth, screenHeight)

	return bg
}

func (bg *Background) ReloadPosition(screenWidth, screenHeight int32) {
	border := bg.GetBorder()
	height := bg.GetHeight()
	posTL := f32Point{float32(border), float32(screenHeight - border)}                          // top-left
	posBL := f32Point{float32(border), float32(screenHeight - (height + border))}               // bottom-left
	posTR := f32Point{float32(screenWidth - border), float32(screenHeight - border)}            // top-right
	posBR := f32Point{float32(screenWidth - border), float32(screenHeight - (height + border))} // bottom-right

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

	bg.gfx.VAOLoad(bg.vao, vertices, gl.STATIC_DRAW)
}

// GetVAO returns the vertex array object associated with the background.
func (bg *Background) GetVAO() *gfx.VAO {
	return bg.vao
}

func (bg *Background) GetBorder() int32 {
	if bg == nil {
		return int32(0)
	}
	return int32(10)
}

func (bg *Background) GetHeight() int32 {
	if bg == nil {
		return int32(0)
	}
	return int32(100)
}
