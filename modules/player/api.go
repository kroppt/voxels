package player

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Interface interface {
	UpdatePlayerPosition(posEvent PositionEvent)
	UpdatePlayerDirection(dirEvent DirectionEvent)
	UpdatePlayerAction(actEvent ActionEvent)
}

// ScrollDirection is either ScrollUp or ScrollDown.
type ScrollDirection int

const (
	ScrollUp   ScrollDirection = 1
	ScrollDown ScrollDirection = 2
)

// ActionEvent contains player action information.
type ActionEvent struct {
	Scroll ScrollDirection
}

// PositionEvent contains player position event information.
//
// X, Y, and Z are voxel coordinates.
type PositionEvent struct {
	X float64
	Y float64
	Z float64
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

// UpdatePlayerAction passes any world altering player actions.
func (m *Module) UpdatePlayerAction(dirEvent ActionEvent) {
	m.c.updateAction(dirEvent)
}

type FnModule struct {
	FnUpdatePlayerPosition  func(posEvent PositionEvent)
	FnUpdatePlayerDirection func(dirEvent DirectionEvent)
	FnUpdatePlayerAction    func(actEvent ActionEvent)
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

func (fn *FnModule) UpdatePlayerAction(actEvent ActionEvent) {
	if fn.FnUpdatePlayerAction != nil {
		fn.FnUpdatePlayerAction(actEvent)
	}
}
