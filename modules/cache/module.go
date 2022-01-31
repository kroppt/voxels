package cache

import (
	"errors"
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
	err := fs.Mkdir("data", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic("failed to create data directory")
	}
	voxelFile, err := fs.OpenFile("data/voxel.data", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("failed to create voxel file")
	}
	chunkFile, err := fs.OpenFile("data/chunk.data", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("failed to create chunk file")
	}
	regionFile, err := fs.OpenFile("data/region.data", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("failed to create region file")
	}
	return &Module{
		c: core{
			voxelFile:    voxelFile,
			chunkFile:    chunkFile,
			regionFile:   regionFile,
			settingsRepo: settingsRepo,
		},
	}
}
