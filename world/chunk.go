package world

import (
	"math/rand"
	"time"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/voxgl"
)

// ChunkPos is a position in chunk space.
type ChunkPos struct {
	X int32
	Z int32
}

// ChunkRange is the range of chunks between Min and Max.
type ChunkRange struct {
	Min ChunkPos
	Max ChunkPos
}

// ForEach executes the given function on every position in the this ChunkRange.
func (rng ChunkRange) ForEach(fn func(pos ChunkPos)) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for z := rng.Min.Z; z <= rng.Max.Z; z++ {
			fn(ChunkPos{X: x, Z: z})
		}
	}
}

// Mul returns this ChunkPos multiplied by another ChunkPos.
func (pos ChunkPos) Mul(s int32) ChunkPos {
	return ChunkPos{
		X: pos.X * s,
		Z: pos.Z * s,
	}
}

// Chunk manages a size X height X size region of voxels.
type Chunk struct {
	Pos      ChunkPos
	flatData []float32
	objs     *voxgl.Object
	root     *Octree
	dirty    bool
	size     int32
	height   int32
}

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size, height int32, pos ChunkPos) *Chunk {
	vertSize := int32(7)
	flatData := make([]float32, size*size*height*vertSize)
	// layout 3+4=7 hard coded in here too
	objs, err := voxgl.NewColoredObject(nil)
	if err != nil {
		panic("failed to make NewColoredObject for chunk")
	}
	chunk := &Chunk{
		Pos:      pos.Mul(size),
		objs:     objs,
		flatData: flatData,
		size:     size,
		height:   height,
	}
	rand.Seed(time.Now().UnixNano())
	for i := int32(0); i < size; i++ {
		for j := int32(0); j < height; j++ {
			for k := int32(0); k < size; k++ {
				x := chunk.Pos.X + i
				y := j
				z := chunk.Pos.Z + k
				r, g, b := rand.Float32(), rand.Float32(), rand.Float32()
				v := Voxel{
					Pos: VoxelPos{x, y, z},
					Col: glm.Vec4{r, g, b, 1.0},
				}
				chunk.SetVoxel(&v)
			}
		}
	}
	return chunk
}

func (c *Chunk) GetRoot() *Octree {
	return c.root
}

// GetRelativeIndices returns the voxel coordinate relative to the origin of
// the chunk, with the assumption that the position is in bounds.
// TODO returns voxel index
func (c *Chunk) GetRelativeIndices(pos VoxelPos) (int, int, int) {
	return int(pos.X - c.Pos.X), int(pos.Y), int(pos.Z - c.Pos.Z)
}

// IsWithinChunk returns whether the position is within the chunk
func (c *Chunk) IsWithinChunk(pos VoxelPos) bool {
	if pos.X < c.Pos.X || pos.Z < c.Pos.Z {
		//pos is below x or z chunk bounds
		return false
	}
	if pos.X >= c.Pos.X+c.size || pos.Z >= c.Pos.Z+c.size {
		// pos is above x or z chunk bounds
		return false
	}
	if pos.Y < 0 || pos.Y >= c.height {
		// y coordinate is out of chunk's bounds
		return false
	}
	return true
}

// SetVoxel updates a voxel's variables in the chunk, if it exists
func (c *Chunk) SetVoxel(v *Voxel) {
	if !c.IsWithinChunk(v.Pos) {
		log.Debugf("%v is not within %v", v, c.Pos)
		return
	}
	x, y, z := float32(v.Pos.X), float32(v.Pos.Y), float32(v.Pos.Z)
	i, j, k := c.GetRelativeIndices(v.Pos)
	r, g, b, a := v.Col[0], v.Col[1], v.Col[2], v.Col[3]
	off := (i + j*int(c.size*c.height) + k*int(c.size)) * 7
	if off%7 != 0 {
		panic("offset not divisible by 7")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	c.flatData[off] = x
	c.flatData[off+1] = y
	c.flatData[off+2] = z
	c.flatData[off+3] = r
	c.flatData[off+4] = g
	c.flatData[off+5] = b
	c.flatData[off+6] = a
	c.root = c.root.AddLeaf(v)
	c.dirty = true
}

// Render renders the chunk in OpenGL.
func (c *Chunk) Render() {
	if c.dirty {
		c.objs.SetData(c.flatData)
		c.dirty = false
	}
	c.objs.Render()
}

func (c *Chunk) Destroy() {
	c.objs.Destroy()
}
