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
	UnloadChunk(chunk.Position)
	UpdateView(map[chunk.Position]struct{}, mgl.Mat4, mgl.Mat4)
	DestroyWindow() error
	Render()
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

// UnloadChunk unloads a chunk.
func (m *Module) UnloadChunk(pos chunk.Position) {
	m.c.unloadChunk(pos)
}

// UpdateView updates what chunks the graphics module should
// try to render.
func (m *Module) UpdateView(viewableChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
	m.c.updateView(viewableChunks, viewMat, projMat)
}

// DestroyWindow destroys an SDL window.
func (m *Module) DestroyWindow() error {
	return m.c.destroyWindow()
}

func (m *Module) Render() {
	m.c.render()
}

type FnModule struct {
	FnCreateWindow  func(string)
	FnShowWindow    func()
	FnPollEvent     func() (sdl.Event, bool)
	FnLoadChunk     func(chunk.Chunk)
	FnUnloadChunk   func(chunk.Position)
	FnUpdateView    func(map[chunk.Position]struct{}, mgl.Mat4, mgl.Mat4)
	FnDestroyWindow func() error
	FnRender        func()
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

func (fn FnModule) UnloadChunk(pos chunk.Position) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(pos)
	}
}

func (fn FnModule) UpdateView(viewableChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
	if fn.FnUpdateView != nil {
		fn.FnUpdateView(viewableChunks, viewMat, projMat)
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
