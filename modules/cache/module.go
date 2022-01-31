package cache

import (
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

type Module struct {
	c core
}

func New(file afero.File, settingsRepo settings.Interface) *Module {
	if settingsRepo == nil {
		panic("cache received nil settings repo")
	}
	return &Module{
		c: core{
			file:         file,
			settingsRepo: settingsRepo,
		},
	}
}
