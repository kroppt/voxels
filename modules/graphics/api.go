package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

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

// UpdateDirection updates the direction of the camera from which to render.
func (m *Module) UpdateDirection(directionEvent DirectionEvent) {
	m.c.updateDirection(directionEvent)
}

// DestroyWindow destroys an SDL window.
func (m *Module) DestroyWindow() error {
	return m.c.destroyWindow()
}
