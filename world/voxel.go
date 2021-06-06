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
func (pos VoxelPos) Add(other VoxelPos) VoxelPos {
	return VoxelPos{
		X: pos.X + other.X,
		Y: pos.Y + other.Y,
		Z: pos.Z + other.Z,
	}
}

// Sub returns this VoxelPos with another VoxelPos subtracted from it.
func (pos VoxelPos) Sub(other VoxelPos) VoxelPos {
	return VoxelPos{
		X: pos.X - other.X,
		Y: pos.Y - other.Y,
		Z: pos.Z - other.Z,
	}
}

// AsVec3 converts this VoxelPos to a glm.Vec3.
func (pos VoxelPos) AsVec3() glm.Vec3 {
	return glm.Vec3{
		float32(pos.X),
		float32(pos.Y),
		float32(pos.Z),
	}
}

// TODO comment
// GetChunkIndex returns the chunk coordinate that the given position
// is in, given the chunkSize.
func (pos VoxelPos) AsChunkPos(chunkSize int32) ChunkPos {
	x := pos.X
	z := pos.Z
	if pos.X < 0 {
		x++
	}
	if pos.Z < 0 {
		z++
	}
	x /= chunkSize
	z /= chunkSize
	if pos.X < 0 {
		x--
	}
	if pos.Z < 0 {
		z--
	}
	return ChunkPos{x, z}
}

// Voxel describes a discrete unit of 3D space.
type Voxel struct {
	Pos VoxelPos
	// TODO Col is RGBA color
	Col glm.Vec4
}
