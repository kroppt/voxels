package world

import (
	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/voxgl"
)

type Voxel struct {
	*voxgl.Object
	coordinates glm.Vec3
}
