package world

import (
	"github.com/engoengine/glm"
)

// VoxelPos is a position in voxel space.
type VoxelPos struct {
	X int32
	Y int32
	Z int32
}

// Add returns this VoxelPos with another VoxelPos added to it.
func (vp VoxelPos) Add(other VoxelPos) VoxelPos {
	return VoxelPos{
		X: vp.X + other.X,
		Y: vp.Y + other.Y,
		Z: vp.Z + other.Z,
	}
}

// Sub returns this VoxelPos with another VoxelPos subtracted from it.
func (vp VoxelPos) Sub(other VoxelPos) VoxelPos {
	return VoxelPos{
		X: vp.X - other.X,
		Y: vp.Y - other.Y,
		Z: vp.Z - other.Z,
	}
}

// AsVec3 converts this VoxelPos to a glm.Vec3.
func (vp VoxelPos) AsVec3() glm.Vec3 {
	return glm.Vec3{
		float32(vp.X),
		float32(vp.Y),
		float32(vp.Z),
	}
}

// Voxel describes a discrete unit of 3D space.
type Voxel struct {
	Pos VoxelPos
	// TODO Col is RGBA color
	Col glm.Vec4
}
