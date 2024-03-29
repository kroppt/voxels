package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/veandco/go-sdl2/sdl"
)

type Interface interface {
	CreateWindow(title string) error
	ShowWindow()
	PollEvent() (sdl.Event, bool)
	LoadChunk(chunk.Chunk)
	UnloadChunk(chunk.ChunkCoordinate)
	UpdateChunk(chunk.Chunk)
	UpdateView(map[chunk.ChunkCoordinate]struct{}, mgl.Mat4)
	UpdateSelection(chunk.VoxelCoordinate, bool)
	DestroyWindow() error
	Render()
	Close()
}

// CreateWindow creates an SDL window.
func (m *Module) CreateWindow(title string) error {
	return m.c.createWindow(title)
}

// ShowWindow makes the current window visible.
func (m *Module) ShowWindow() {
	m.c.showWindow()
}

// PollEvent returns the next event if present and whether it was present.
func (m *Module) PollEvent() (sdl.Event, bool) {
	return m.c.pollEvent()
}

// LoadChunk loads a chunk.
func (m *Module) LoadChunk(chunk chunk.Chunk) {
	m.c.loadChunk(chunk)
}

// UpdateChunk updates a chunk.
func (m *Module) UpdateChunk(chunk chunk.Chunk) {
	m.c.updateChunk(chunk)
}

// UnloadChunk unloads a chunk.
func (m *Module) UnloadChunk(pos chunk.ChunkCoordinate) {
	m.c.unloadChunk(pos)
}

// UpdateView updates what chunks the graphics module should
// try to render.
func (m *Module) UpdateView(viewableChunks map[chunk.ChunkCoordinate]struct{}, viewMat mgl.Mat4) {
	m.c.updateView(viewableChunks, viewMat)
}

// UpdateSelection updates the currently selected voxel
func (m *Module) UpdateSelection(selectedVoxel chunk.VoxelCoordinate, selected bool) {
	m.c.updateSelection(selectedVoxel, selected)
}

// DestroyWindow destroys an SDL window.
func (m *Module) DestroyWindow() error {
	return m.c.destroyWindow()
}

func (m *Module) Render() {
	m.c.render()
}

// Close does nothing.
func (m *Module) Close() {
}

type FnModule struct {
	FnCreateWindow    func(string)
	FnShowWindow      func()
	FnPollEvent       func() (sdl.Event, bool)
	FnLoadChunk       func(chunk.Chunk)
	FnUpdateChunk     func(chunk.Chunk)
	FnUnloadChunk     func(chunk.ChunkCoordinate)
	FnUpdateView      func(map[chunk.ChunkCoordinate]struct{}, mgl.Mat4)
	FnUpdateSelection func(chunk.VoxelCoordinate, bool)
	FnDestroyWindow   func() error
	FnRender          func()
	FnClose           func()
}

func (fn FnModule) CreateWindow(title string) error {
	if fn.FnCreateWindow != nil {
		fn.FnCreateWindow(title)
	}
	return nil
}

func (fn FnModule) ShowWindow() {
	if fn.FnShowWindow != nil {
		fn.FnShowWindow()
	}
}

func (fn FnModule) PollEvent() (sdl.Event, bool) {
	if fn.FnPollEvent != nil {
		return fn.FnPollEvent()
	}
	return nil, false
}

func (fn FnModule) LoadChunk(chunk chunk.Chunk) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(chunk)
	}
}

func (fn FnModule) UpdateChunk(chunk chunk.Chunk) {
	if fn.FnUpdateChunk != nil {
		fn.FnUpdateChunk(chunk)
	}
}

func (fn FnModule) UnloadChunk(pos chunk.ChunkCoordinate) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(pos)
	}
}

func (fn FnModule) UpdateView(viewableChunks map[chunk.ChunkCoordinate]struct{}, viewMat mgl.Mat4) {
	if fn.FnUpdateView != nil {
		fn.FnUpdateView(viewableChunks, viewMat)
	}
}

func (fn FnModule) UpdateSelection(selectedVoxel chunk.VoxelCoordinate, selected bool) {
	if fn.FnUpdateSelection != nil {
		fn.FnUpdateSelection(selectedVoxel, selected)
	}
}

func (fn FnModule) DestroyWindow() error {
	if fn.FnDestroyWindow != nil {
		return fn.FnDestroyWindow()
	}
	return nil
}

func (fn FnModule) Render() {
	if fn.FnRender != nil {
		fn.FnRender()
	}
}

func (fn FnModule) Close() {
	if fn.FnClose != nil {
		fn.FnClose()
	}
}
