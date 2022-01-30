package tick

import "github.com/kroppt/voxels/modules/camera"

type Module struct {
	c core
}

func New(cameraMod camera.Interface, timeMod Time, tickRateNano int64) *Module {
	if cameraMod == nil {
		panic("camera module was nil")
	}
	if timeMod == nil {
		panic("time module was nil")
	}
	if tickRateNano <= 0 {
		panic("tick rate was not positive")
	}
	return &Module{
		core{
			cameraMod:    cameraMod,
			timeMod:      timeMod,
			lastTickNano: timeMod.Now().UnixNano(),
			tickRateNano: tickRateNano,
			currentTick:  0,
		},
	}
}
