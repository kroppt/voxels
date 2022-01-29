package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/veandco/go-sdl2/sdl"
)

type Interface interface {
	CreateWindow(title string, width, height uint32) error
	ShowWindow()
	PollEvent() (sdl.Event, bool)
	LoadChunk(chunk.Chunk)
	UnloadChunk(chunk.Position)
	UpdateViewableChunks(map[chunk.Position]struct{})
	DestroyWindow() error
	Render()
}

// DirectionEvent contains rotation information.
type DirectionEvent struct {
	Rotation mgl.Quat
}

// CreateWindow creates an SDL window.
func (m *Module) CreateWindow(title string, width, height uint32) error {
	return m.c.createWindow(title, width, height)
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
}

// UnloadChunk unloads a chunk.
func (m *Module) UnloadChunk(chunk.Position) {
}

// UpdateViewableChunks updates what chunks the graphics module should
// try to render.
func (m *Module) UpdateViewableChunks(map[chunk.Position]struct{}) {

}

// DestroyWindow destroys an SDL window.
func (m *Module) DestroyWindow() error {
	return m.c.destroyWindow()
}

func (m *Module) Render() {
	m.c.render()
}

type FnModule struct {
	FnCreateWindow         func(string, uint32, uint32)
	FnShowWindow           func()
	FnPollEvent            func() (sdl.Event, bool)
	FnLoadChunk            func(chunk.Chunk)
	FnUnloadChunk          func(chunk.Position)
	FnUpdateViewableChunks func(map[chunk.Position]struct{})
	FnDestroyWindow        func() error
	FnRender               func()
}

func (fn FnModule) CreateWindow(title string, width, height uint32) error {
	if fn.FnCreateWindow != nil {
		fn.FnCreateWindow(title, width, height)
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

func (fn FnModule) UpdateViewableChunks(viewableChunks map[chunk.Position]struct{}) {
	if fn.FnUpdateViewableChunks != nil {
		fn.FnUpdateViewableChunks(viewableChunks)
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
