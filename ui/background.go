package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*Background)(nil)

// Background stores the internal data associated with a background UI element.
type Background struct {
	vao *gfx.VAO
}

// NewBackground creates a Background element.
func NewBackground(gfx Gfx) *Background {
	layout := []int32{2, 4}
	posTL := f32Point{-1, 1}   // top-left
	posBL := f32Point{-1, 0.5} // bottom-left
	posTR := f32Point{1, 1}    // top-right
	posBR := f32Point{1, 0.5}  // bottom-right

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

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	gfx.VAOLoad(vao, vertices, gl.STATIC_DRAW)

	bg := &Background{
		vao: vao,
	}
	return bg
}

// GetVAO returns the vertex array object associated with the background.
func (bg *Background) GetVAO() *gfx.VAO {
	return bg.vao
}
