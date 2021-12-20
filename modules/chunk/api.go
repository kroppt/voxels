package chunk

type Interface interface {
	UpdatePosition(posEvent PositionEvent)
}

// PositionEvent contains position event information.
type PositionEvent struct {
	X int32
	Y int32
	Z int32
}

// UpdatePosition updates the chunks based on the new position.
func (m *Module) UpdatePosition(posEvent PositionEvent) {
	m.c.updatePosition(posEvent)
}

type FnModule struct {
	FnUpdatePosition func(posEvent PositionEvent)
}

func (fn *FnModule) UpdatePosition(posEvent PositionEvent) {
	if fn.FnUpdatePosition != nil {
		fn.FnUpdatePosition(posEvent)
	}
}
