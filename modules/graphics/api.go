package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

type Interface interface {
	CreateWindow(title string, width, height uint32) error
	ShowWindow()
	PollEvent() (sdl.Event, bool)
	UpdatePlayerDirection(directionEvent DirectionEvent)
	ShowChunk(chunkEvent ChunkEvent)
	DestroyWindow() error
}

// DirectionEvent contains rotation information.
type DirectionEvent struct {
	Rotation mgl.Quat
}

// ChunkEvent contains chunk information.
type ChunkEvent struct {
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

// UpdatePlayerDirection updates the direction of the camera from which to render.
func (m *Module) UpdatePlayerDirection(directionEvent DirectionEvent) {
	m.c.updatePlayerDirection(directionEvent)
}

// ShowChunk shows a chunk.
func (m *Module) ShowChunk(chunkEvent ChunkEvent) {
	m.c.showChunk(chunkEvent)
}

// DestroyWindow destroys an SDL window.
func (m *Module) DestroyWindow() error {
	return m.c.destroyWindow()
}

type FnModule struct {
	FnCreateWindow          func(string, uint32, uint32)
	FnShowWindow            func()
	FnPollEvent             func() (sdl.Event, bool)
	FnUpdatePlayerDirection func(DirectionEvent)
	FnShowChunk             func(chunkEvent ChunkEvent)
	FnDestroyWindow         func() error
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

func (fn FnModule) UpdatePlayerDirection(directionEvent DirectionEvent) {
	if fn.FnUpdatePlayerDirection != nil {
		fn.FnUpdatePlayerDirection(directionEvent)
	}
}

func (fn FnModule) ShowChunk(chunkEvent ChunkEvent) {
	if fn.FnShowChunk != nil {
		fn.FnShowChunk(chunkEvent)
	}
}

func (fn FnModule) DestroyWindow() error {
	if fn.FnDestroyWindow != nil {
		return fn.FnDestroyWindow()
	}
	return nil
}
