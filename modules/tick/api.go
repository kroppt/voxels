package tick

type Interface interface {
	GetTick() int
	AdvanceTick()
}

func (m *Module) GetTick() int {
	return m.c.getTick()
}

func (m *Module) AdvanceTick() {
	m.c.advanceTick()
}

type FnModule struct {
	FnGetTick     func() int
	FnAdvanceTick func()
}

func (fn FnModule) GetTick() int {
	if fn.FnGetTick != nil {
		return fn.FnGetTick()
	}
	return -1
}

func (fn FnModule) AdvanceTick() {
	if fn.FnAdvanceTick != nil {
		fn.FnAdvanceTick()
	}
}
