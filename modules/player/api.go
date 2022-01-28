package player

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Interface interface {
	UpdatePlayerPosition(posEvent PositionEvent)
	UpdatePlayerDirection(dirEvent DirectionEvent)
}

// PositionEvent contains player position event information.
//
// X, Y, and Z are voxel coordinates.
type PositionEvent struct {
	X int32
	Y int32
	Z int32
}

// DirectionEvent contains rotation information.
type DirectionEvent struct {
	Rotation mgl.Quat
}

// UpdatePlayerPosition updates the chunks based on the new player position.
func (m *Module) UpdatePlayerPosition(posEvent PositionEvent) {
	m.c.updatePosition(posEvent)
}

// UpdatePlayerDirection updates what chunks the player should see based on look direction.
func (m *Module) UpdatePlayerDirection(dirEvent DirectionEvent) {
	m.c.updateDirection(dirEvent)
}

type FnModule struct {
	FnUpdatePlayerPosition  func(posEvent PositionEvent)
	FnUpdatePlayerDirection func(dirEvent DirectionEvent)
}

func (fn *FnModule) UpdatePlayerPosition(posEvent PositionEvent) {
	if fn.FnUpdatePlayerPosition != nil {
		fn.FnUpdatePlayerPosition(posEvent)
	}
}

func (fn *FnModule) UpdatePlayerDirection(dirEvent DirectionEvent) {
	if fn.FnUpdatePlayerDirection != nil {
		fn.FnUpdatePlayerDirection(dirEvent)
	}
}
