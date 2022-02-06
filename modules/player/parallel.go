package player

// UpdatePlayerPosition updates the chunks based on the new player position.
func (m *ParallelModule) UpdatePlayerPosition(posEvent PositionEvent) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updatePosition(posEvent)
		close(done)
	}
	<-done
}

// UpdatePlayerDirection updates what chunks the player should see based on look direction.
func (m *ParallelModule) UpdatePlayerDirection(dirEvent DirectionEvent) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updateDirection(dirEvent)
		close(done)
	}
	<-done
}

// UpdatePlayerAction passes any world altering player actions.
func (m *ParallelModule) UpdatePlayerAction(dirEvent ActionEvent) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updateAction(dirEvent)
		close(done)
	}
	<-done
}
