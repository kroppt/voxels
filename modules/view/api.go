package view

import (
	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/chunk"
)

type Interface interface {
	AddTree(chunk.ChunkCoordinate, *Octree)
	RemoveTree(chunk.ChunkCoordinate)
	AddNode(chunk.VoxelCoordinate)
	RemoveNode(chunk.VoxelCoordinate)
	UpdateView(ViewState)
	GetSelection() (chunk.VoxelCoordinate, bool)
	GetPlacement() (chunk.VoxelCoordinate, bool)
}

// AddTree adds an entire tree to the map of trees by chunk coord.
func (m *Module) AddTree(cc chunk.ChunkCoordinate, root *Octree) {
	m.c.addTree(cc, root)
}

// RemoveTree removes an entire tree from the map of trees.
func (m *Module) RemoveTree(cc chunk.ChunkCoordinate) {
	m.c.removeTree(cc)
}

// AddNode adds a voxel to its chunk's tree, panics if no tree found.
func (m *Module) AddNode(vc chunk.VoxelCoordinate) {
	m.c.addNode(vc)
}

// RemoveNode removes a voxel from its chunk's tree, panics if no tree found.
func (m *Module) RemoveNode(vc chunk.VoxelCoordinate) {
	m.c.removeNode(vc)
}

// UpdateView calculated frustum culling and the view matrix and passes it to
// the graphics module.
func (m *Module) UpdateView(vs ViewState) {
	m.c.updateView(vs)
}

// GetSelection re-calculates the currently selected voxel and returns it. The return
// value should only be considered if it is paired with a true boolean return.
func (m *Module) GetSelection() (chunk.VoxelCoordinate, bool) {
	return m.c.getSelection()
}

// GetPlacement re-calculates the current placement voxel and returns it. The return
// value should only be considered if it is paired with a true boolean return.
// i.e. this is where you would place a block if you were to place one
func (m *Module) GetPlacement() (chunk.VoxelCoordinate, bool) {
	return m.c.getPlacement()
}

// ViewState represents the way the view is positioned and oriented.
type ViewState struct {
	Pos mgl.Vec3
	Dir mgl.Quat
}

type FnModule struct {
	FnAddTree      func(chunk.ChunkCoordinate, *Octree)
	FnRemoveTree   func(chunk.ChunkCoordinate)
	FnAddNode      func(chunk.VoxelCoordinate)
	FnRemoveNode   func(chunk.VoxelCoordinate)
	FnUpdateView   func(ViewState)
	FnGetSelection func() (chunk.VoxelCoordinate, bool)
	FnGetPlacement func() (chunk.VoxelCoordinate, bool)
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

func (fn *FnModule) GetSelection() (chunk.VoxelCoordinate, bool) {
	if fn.FnGetSelection != nil {
		return fn.FnGetSelection()
	}
	return chunk.VoxelCoordinate{}, false
}

func (fn *FnModule) GetPlacement() (chunk.VoxelCoordinate, bool) {
	if fn.FnGetPlacement != nil {
		return fn.FnGetPlacement()
	}
	return chunk.VoxelCoordinate{}, false
}

func (fn *FnModule) AddTree(cc chunk.ChunkCoordinate, root *Octree) {
	if fn.FnAddTree != nil {
		fn.FnAddTree(cc, root)
	}
}

func (fn *FnModule) RemoveTree(cc chunk.ChunkCoordinate) {
	if fn.FnRemoveTree != nil {
		fn.FnRemoveTree(cc)
	}
}
