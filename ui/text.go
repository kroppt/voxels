package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/log"
)

var _ Element = (*Text)(nil)

type Text struct {
	value   string
	program *gfx.Program
	vao     *gfx.VAO
	gfx     Gfx
	fnt     *gfx.FontInfo
	parent  *Button
}

func NewText(gfx Gfx, parentComponent *Button, screenWidth, screenHeight int32, value string) (*Text, error) {
	layout := []int32{2, 2}

	vao := gfx.NewVAO(gl.TRIANGLES, layout)

	fnt, err := gfx.LoadFontTexture("NotoMono-Regular.ttf", fontSize)
	if err != nil {
		return nil, err
	}

	text := &Text{
		value:  value,
		vao:    vao,
		gfx:    gfx,
		fnt:    fnt,
		parent: parentComponent,
	}

	text.ReloadPosition(screenWidth, screenHeight)

	return text, nil
}

func (text *Text) GetProgram() *gfx.Program {
	return text.program
}

func (text *Text) SetProgram(program *gfx.Program) {
	text.program = program
}

func (text *Text) GetVAO() *gfx.VAO {
	return text.vao
}

func (text *Text) FontTextureBind() {
	text.fnt.GetTexture().Bind()
}

func (text *Text) FontTextureUnbind() {
	text.fnt.GetTexture().Unbind()
}

func (text *Text) GetBorder() int32 {
	if text == nil {
		return int32(0)
	}
	return int32(10)
}

func (text *Text) ReloadPosition(screenWidth, screenHeight int32) {
	pos := gfx.Point{X: text.parent.GetLeft() + text.GetBorder(), Y: screenHeight - text.parent.GetTop() - text.GetBorder()}
	align := gfx.Align{V: gfx.AlignBelow, H: gfx.AlignLeft}
	textTriangles := text.fnt.MapString(text.value, pos, align)

	err := text.gfx.VAOLoad(text.vao, textTriangles, gl.STATIC_DRAW)

	if err != nil {
		log.Fatalf("failed to load button text triangles: %v", err)
	}
}
