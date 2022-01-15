package chunk

type Interface interface {
	UpdatePlayerPosition(posEvent PositionEvent)
}

// PositionEvent contains player position event information.
//
// X, Y, and Z are voxel coordinates.
type PositionEvent struct {
	X int32
	Y int32
	Z int32
}

// UpdatePlayerPosition updates the chunks based on the new player position.
func (m *Module) UpdatePlayerPosition(posEvent PositionEvent) {
	m.c.updatePosition(posEvent)
}

type FnModule struct {
	FnUpdatePlayerPosition func(posEvent PositionEvent)
}

func (fn *FnModule) UpdatePlayerPosition(posEvent PositionEvent) {
	if fn.FnUpdatePlayerPosition != nil {
		fn.FnUpdatePlayerPosition(posEvent)
	}
}
