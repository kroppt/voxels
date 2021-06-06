// +build test

package voxgl

import (
	"github.com/kroppt/gfx"
)

func NewColoredObject(vertices []float32) (*Object, error) {
	obj, _ := NewObject(gfx.Program{}, []float32{}, []int32{})
	return obj, nil
}
