package world

import (
	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/voxgl"
)

type Voxel struct {
	*voxgl.Object
	Pos glm.Vec3
	Col glm.Vec4
}
