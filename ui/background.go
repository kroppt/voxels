package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*Background)(nil)

// Background stores the internal data associated with a background UI element.
type Background struct {
	vao   *gfx.VAO
	color [4]float32
}

// NewBackground creates a Background element.
func NewBackground(gfx Gfx) *Background {
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

	var color [4]float32
	color[0] = 1.0 // red
	color[1] = 0.0 // green
	color[2] = 0.0 // blue
	color[3] = 1.0 // alpha

	bg := &Background{
		vao:   vao,
		color: color,
	}
	return bg
}

// GetVAO returns the vertex array object associated with the background.
func (bg *Background) GetVAO() *gfx.VAO {
	return bg.vao
}

// GetColor returns a color array of the RGBA of the background.
func (bg *Background) GetColor() *[4]float32 {
	return &bg.color
}
