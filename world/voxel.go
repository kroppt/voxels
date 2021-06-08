package world

import (
	"github.com/engoengine/glm"
)

// VoxelPos is a position in voxel space.
type VoxelPos struct {
	X int
	Y int
	Z int
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

// GetChunkPos returns the chunk position that the this VoxelPos is in for the
// given chunkSize.
func (pos VoxelPos) GetChunkPos(chunkSize int) ChunkPos {
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

// AsLocalChunkPos returns the voxel position relative to the origin of chunk,
// with the assumption that the position is in the bounds of the chunk.
func (pos VoxelPos) AsLocalChunkPos(chunk Chunk) VoxelPos {
	return VoxelPos{
		X: pos.X - chunk.AsVoxelPos().X,
		Y: pos.Y, // TODO this will break with 3D chunks
		Z: pos.Z - chunk.AsVoxelPos().Z,
	}
}

// Color is an RGBA color.
type Color struct {
	R float32
	G float32
	B float32
	A float32
}

// Voxel describes a discrete unit of 3D space.
type Voxel struct {
	Pos   VoxelPos
	Color Color
}
