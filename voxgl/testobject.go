// +build test

package voxgl

import (
	"github.com/kroppt/gfx"
)

type Object struct {
}

func NewObject(program gfx.Program, vertices []float32, layout []int32) (*Object, error) {
	return &Object{}, nil
}

func (o *Object) SetData(data []float32) {
}

func (o *Object) Render() {
}

func (o *Object) Translate(x, y, z float32) {
}

func (o *Object) Scale(x, y, z float32) {
}

func (o *Object) Rotate(x, y, z float32) {
}

func (o *Object) Destroy() {
}
