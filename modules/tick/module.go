package tick

import "github.com/kroppt/voxels/modules/camera"

type Module struct {
	c core
}

func New(cameraMod camera.Interface) *Module {
	if cameraMod == nil {
		panic("camera module was nil")
	}
	return &Module{
		core{
			cameraMod:   cameraMod,
			currentTick: 0,
		},
	}
}
