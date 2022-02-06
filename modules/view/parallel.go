package view

import "github.com/kroppt/voxels/chunk"

// AddTree adds an entire tree to the map of trees by chunk coord.
func (m *ParallelModule) AddTree(cc chunk.ChunkCoordinate, root *Octree) {
	m.do <- func() { m.c.addTree(cc, root) }
}

// RemoveTree removes an entire tree from the map of trees.
func (m *ParallelModule) RemoveTree(cc chunk.ChunkCoordinate) {
	m.do <- func() { m.c.removeTree(cc) }
}

// AddNode adds a voxel to its chunk's tree, panics if no tree found.
func (m *ParallelModule) AddNode(vc chunk.VoxelCoordinate) {
	m.do <- func() { m.c.addNode(vc) }
}

// RemoveNode removes a voxel from its chunk's tree, panics if no tree found.
func (m *ParallelModule) RemoveNode(vc chunk.VoxelCoordinate) {
	m.do <- func() { m.c.removeNode(vc) }
}

// UpdateView calculated frustum culling and the view matrix and passes it to
// the graphics module.
func (m *ParallelModule) UpdateView(vs ViewState) {
	m.do <- func() { m.c.updateView(vs) }
}

// GetSelection re-calculates the currently selected voxel and returns it. The return
// value should only be considered if it is paired with a true boolean return.
func (m *ParallelModule) GetSelection() (chunk.VoxelCoordinate, bool) {
	type returns struct {
		vc chunk.VoxelCoordinate
		ok bool
	}
	done := make(chan returns)
	m.do <- func() {
		vc, ok := m.c.getSelection()
		done <- returns{vc, ok}
	}
	d := <-done
	return d.vc, d.ok
}

// GetPlacement re-calculates the current placement voxel and returns it. The return
// value should only be considered if it is paired with a true boolean return.
// i.e. this is where you would place a block if you were to place one
func (m *ParallelModule) GetPlacement() (chunk.VoxelCoordinate, bool) {
	type returns struct {
		vc chunk.VoxelCoordinate
		ok bool
	}
	done := make(chan returns)
	m.do <- func() {
		vc, ok := m.c.getPlacement()
		done <- returns{vc, ok}
	}
	d := <-done
	return d.vc, d.ok
}
