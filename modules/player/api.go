package player

// MoveDirection represents a direction of movement.
type MoveDirection int

const (
	// MoveForwards represents moving forward.
	MoveForwards MoveDirection = 1
	// MoveRight represents strafing right.
	MoveRight MoveDirection = 2
	// MoveBackwards represents moving back.
	MoveBackwards MoveDirection = 3
	// MoveLeft represents strafing left.
	MoveLeft MoveDirection = 4
)

// MovementEvent contains movement event information.
type MovementEvent struct {
	Direction MoveDirection
}

// LookEvent contains look event information.
type LookEvent struct {
	Right int32
	Down  int32
}

// HandleMovementEvent handles a movement event.
func (m *Module) HandleMovementEvent(evt MovementEvent) {
}

// HandleLookEvent handles a look event.
func (m *Module) HandleLookEvent(evt LookEvent) {
}
