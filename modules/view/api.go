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

// AddNode adds a voxel to the underlying octree, to be considered for selection.
func (m *Module) AddNode(vc chunk.VoxelCoordinate) {
	m.c.addNode(vc)
}

// RemoveNode removes a voxel from the underlying octree, to no longer be considered for selection.
func (m *Module) RemoveNode(vc chunk.VoxelCoordinate) {
	m.c.removeNode(vc)
}

// UpdateView calculated frustum culling and the view matrix and passes it to
// the graphics module.
func (m *Module) UpdateView(vs ViewState) {
	m.c.updateView(vs)
}

// UpdateSelection re-calculates the currently selected voxel and passes it to
// the graphics module.
func (m *Module) UpdateSelection() {
	m.c.updateSelection()
}

// GetSelection re-calculates the currently selected voxel and returns it. The return
// value should only be considered if it is paired with a true boolean return.
func (m *Module) GetSelection() (chunk.VoxelCoordinate, bool) {
	return m.c.getSelection()
}

// ViewState represents the way the view is positioned and oriented.
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
