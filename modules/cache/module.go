package cache

import (
	"os"

	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

type Module struct {
	c core
}

func New(fs afero.Fs, settingsRepo settings.Interface) *Module {
	if settingsRepo == nil {
		panic("cache received nil settings repo")
	}
	f, err := fs.OpenFile("world.data", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("failed to create file")
	}
	return &Module{
		c: core{
			file:         f,
			settingsRepo: settingsRepo,
		},
	}
}
