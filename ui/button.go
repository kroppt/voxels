package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

var _ Element = (*Button)(nil)

// Button stores the internal data associated with a button UI element.
type Button struct {
	vao *gfx.VAO
}

func NewButton(gfx Gfx) *Button {
	layout := []int32{2, 4}
	border := float32(0.125)
	posTL := f32Point{-1 + border, 1 - border}     // top-left
	posBL := f32Point{-1 + border, 0.5 + border}   // bottom-left
	posTR := f32Point{-0.5 - border, 1 - border}   // top-right
	posBR := f32Point{-0.5 - border, 0.5 + border} // bottom-right

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

	vao := gfx.NewVAO(gl.TRIANGLE_STRIP, layout)

	gfx.VAOLoad(vao, vertices, gl.STATIC_DRAW)

	btn := &Button{
		vao: vao,
	}
	return btn
}

// GetVAO returns the vertex array object associated with the button.
func (btn *Button) GetVAO() *gfx.VAO {
	return btn.vao
}
