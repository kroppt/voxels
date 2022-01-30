package camera

type Interface interface {
	HandleMovementEvent(MovementEvent)
	HandleLookEvent(LookEvent)
	Tick()
}

// MoveDirection represents a direction of movement.
type MoveDirection int

func (d MoveDirection) String() string {
	switch d {
	case MoveForwards:
		return "forwards"
	case MoveRight:
		return "right"
	case MoveBackwards:
		return "backwards"
	case MoveLeft:
		return "left"
	}
	return "invalid"
}

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
	Pressed   bool
}

// LookEvent contains look event information.
type LookEvent struct {
	Right float64
	Down  float64
}

func (m *Module) Tick() {
	m.c.tick()
}

// HandleMovementEvent handles a movement event.
func (m *Module) HandleMovementEvent(evt MovementEvent) {
	m.c.handleMovementEvent(evt)
}

// HandleLookEvent handles a look event.
func (m *Module) HandleLookEvent(evt LookEvent) {
	m.c.handleLookEvent(evt)
}

type FnModule struct {
	FnHandleMovementEvent func(MovementEvent)
	FnHandleLookEvent     func(LookEvent)
	FnTick                func()
}

func (fn *FnModule) HandleMovementEvent(movementEvent MovementEvent) {
	if fn.FnHandleMovementEvent != nil {
		fn.FnHandleMovementEvent(movementEvent)
	}
}

func (fn *FnModule) HandleLookEvent(lookEvent LookEvent) {
	if fn.FnHandleLookEvent != nil {
		fn.FnHandleLookEvent(lookEvent)
	}
}

func (fn *FnModule) Tick() {
	if fn.FnTick != nil {
		fn.FnTick()
	}
}
