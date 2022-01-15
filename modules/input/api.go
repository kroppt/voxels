package input

type Interface interface {
	RouteEvents()
}

// RouteEvents polls for input events and routes them to other modules.
func (m *Module) RouteEvents() {
	m.c.routeEvents()
}

// PixelsToRadians converts from pixels to radians in terms of camera rotation.
func (m *Module) PixelsToRadians(xRel, yRel int32) (float64, float64) {
	return m.c.pixelsToRadians(xRel, yRel)
}

type FnModule struct {
	FnRouteEvents     func()
	FnPixelsToRadians func(xRel, yRel int32) (float64, float64)
}

func (fn *FnModule) RouteEvents() {
	if fn.FnRouteEvents != nil {
		fn.FnRouteEvents()
	}
}

func (fn *FnModule) PixelsToRadians(xRel, yRel int32) (float64, float64) {
	if fn.FnPixelsToRadians != nil {
		return fn.FnPixelsToRadians(xRel, yRel)
	}
	return 0, 0
}
