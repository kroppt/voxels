package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

type Interface interface {
	CreateWindow(title string, width, height uint32) error
	ShowWindow()
	PollEvent() (sdl.Event, bool)
	LoadChunk(ChunkEvent)
	UnloadChunk(ChunkEvent)
	UpdateViewableChunks(map[ChunkEvent]struct{})
	DestroyWindow() error
	Render()
}

// PositionEvent contains position information.
//
// X, Y, and Z are sub-voxel coordinates.
type PositionEvent struct {
	X float64
	Y float64
	Z float64
}

// DirectionEvent contains rotation information.
type DirectionEvent struct {
	Rotation mgl.Quat
}

// ChunkEvent contains chunk information.
type ChunkEvent struct {
	PositionX int32
	PositionY int32
	PositionZ int32
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
func (m *Module) LoadChunk(chunkEvent ChunkEvent) {
}

// UnloadChunk unloads a chunk.
func (m *Module) UnloadChunk(chunkEvent ChunkEvent) {
}

// UpdateViewableChunks updates what chunks the graphics module should
// try to render.
func (m *Module) UpdateViewableChunks(map[ChunkEvent]struct{}) {

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
	FnLoadChunk            func(chunkEvent ChunkEvent)
	FnUnloadChunk          func(chunkEvent ChunkEvent)
	FnUpdateViewableChunks func(map[ChunkEvent]struct{})
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

func (fn FnModule) LoadChunk(chunkEvent ChunkEvent) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(chunkEvent)
	}
}

func (fn FnModule) UnloadChunk(chunkEvent ChunkEvent) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(chunkEvent)
	}
}

func (fn FnModule) UpdateViewableChunks(viewableChunks map[ChunkEvent]struct{}) {
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
