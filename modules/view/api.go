package view

import (
	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/chunk"
)

type Interface interface {
	AddNode(chunk.VoxelCoordinate)
	RemoveNode(chunk.VoxelCoordinate)
	UpdateView(ViewState)
	UpdateSelection()
	GetSelection() (chunk.VoxelCoordinate, bool)
}

func (m *Module) AddNode(vc chunk.VoxelCoordinate) {
	m.c.addNode(vc)
}

func (m *Module) RemoveNode(vc chunk.VoxelCoordinate) {
	m.c.removeNode(vc)
}

func (m *Module) UpdateView(vs ViewState) {
	m.c.updateView(vs)
}

func (m *Module) UpdateSelection() {
	m.c.updateSelection()
}

func (m *Module) GetSelection() (chunk.VoxelCoordinate, bool) {
	return m.c.getSelection()
}

type ViewState struct {
	Pos mgl.Vec3
	Dir mgl.Quat
}

type FnModule struct {
	FnAddNode         func(chunk.VoxelCoordinate)
	FnRemoveNode      func(chunk.VoxelCoordinate)
	FnUpdateView      func(ViewState)
	FnUpdateSelection func()
	FnGetSelection    func() (chunk.VoxelCoordinate, bool)
}

func (fn *FnModule) AddNode(vc chunk.VoxelCoordinate) {
	if fn.FnAddNode != nil {
		fn.FnAddNode(vc)
	}
}
func (fn *FnModule) RemoveNode(vc chunk.VoxelCoordinate) {
	if fn.FnRemoveNode != nil {
		fn.FnRemoveNode(vc)
	}
}
func (fn *FnModule) UpdateView(vs ViewState) {
	if fn.FnUpdateView != nil {
		fn.FnUpdateView(vs)
	}
}
func (fn *FnModule) UpdateSelection() {
	if fn.FnUpdateSelection != nil {
		fn.FnUpdateSelection()
	}
}
func (fn *FnModule) GetSelection() (chunk.VoxelCoordinate, bool) {
	if fn.FnGetSelection != nil {
		return fn.FnGetSelection()
	}
	return chunk.VoxelCoordinate{}, false
}
