package tick

type Interface interface {
	GetTick() int
	AdvanceTick()
	IsNextTickReady() bool
}

func (m *Module) GetTick() int {
	return m.c.getTick()
}

func (m *Module) AdvanceTick() {
	m.c.advanceTick()
}

func (m *Module) IsNextTickReady() bool {
	return m.c.isNextTickReady()
}

type FnModule struct {
	FnGetTick         func() int
	FnAdvanceTick     func()
	FnIsNextTickReady func() bool
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

func (fn FnModule) IsNextTickReady() bool {
	if fn.FnIsNextTickReady != nil {
		return fn.FnIsNextTickReady()
	}
	return false
}
