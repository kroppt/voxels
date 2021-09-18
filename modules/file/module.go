package file

// Module is a file system API.
type Module struct {
	c core
}

// New creates a file module.
func New() *Module {
	return &Module{}
}
