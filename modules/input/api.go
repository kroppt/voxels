package input

// RouteEvents polls for input events and routes them to other modules.
func (m *Module) RouteEvents() {
	m.c.routeEvents()
}

// PixelsToRadians converts from pixels to radians in terms of camera rotation.
func (m *Module) PixelsToRadians(xRel, yRel int32) (float64, float64) {
	return m.c.pixelsToRadians(xRel, yRel)
}
