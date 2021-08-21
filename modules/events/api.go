package events

// RouteEvents polls for SDL events and routes them to other modules.
func (m *Module) RouteEvents() {
	m.c.routeEvents()
}
