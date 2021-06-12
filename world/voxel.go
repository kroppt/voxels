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

// VoxelRange is the range of voxels between Min and Max.
type VoxelRange struct {
	Min VoxelPos
	Max VoxelPos
}

// GetSurroundings returns a range surrounding the position by amount in every
// direction.
func (pos VoxelPos) GetSurroundings(amount int) VoxelRange {
	minx := pos.X - amount
	maxx := pos.X + amount
	miny := pos.Y - amount
	maxy := pos.Y + amount
	mink := pos.Z - amount
	maxk := pos.Z + amount
	return VoxelRange{
		Min: VoxelPos{minx, miny, mink},
		Max: VoxelPos{maxx, maxy, maxk},
	}
}

// ForEach executes the given function on every position in the this VoxelRange.
func (rng VoxelRange) ForEach(fn func(pos VoxelPos)) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for y := rng.Min.Y; y <= rng.Max.Y; y++ {
			for z := rng.Min.Z; z <= rng.Max.Z; z++ {
				fn(VoxelPos{X: x, Y: y, Z: z})
			}
		}
	}
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
	y := pos.Y
	z := pos.Z
	if pos.X < 0 {
		x++
	}
	if pos.Y < 0 {
		y++
	}
	if pos.Z < 0 {
		z++
	}
	x /= chunkSize
	y /= chunkSize
	z /= chunkSize
	if pos.X < 0 {
		x--
	}
	if pos.Y < 0 {
		y--
	}
	if pos.Z < 0 {
		z--
	}
	return ChunkPos{x, y, z}
}

// AsLocalChunkPos returns the voxel position relative to the origin of chunk,
// with the assumption that the position is in the bounds of the chunk.
func (pos VoxelPos) AsLocalChunkPos(chunk Chunk) VoxelPos {
	return VoxelPos{
		X: pos.X - chunk.AsVoxelPos().X,
		Y: pos.Y - chunk.AsVoxelPos().Y,
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
	Pos     VoxelPos
	Color   Color
	AdjMask AdjacentMask
	Btype   BlockType
}
