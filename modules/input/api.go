package input

// RouteEvents polls for input events and routes them to other modules.
func (m *Module) RouteEvents() {
	m.c.routeEvents()
}
