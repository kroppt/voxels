package chunk

type Interface interface {
	UpdatePlayerPosition(posEvent PositionEvent)
}

// PositionEvent contains player position event information.
type PositionEvent struct {
	X float64
	Y float64
	Z float64
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
