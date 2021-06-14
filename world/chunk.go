package world

import (
	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/voxgl"
)

// ChunkPos is a position in chunk space.
type ChunkPos struct {
	X int
	Y int
	Z int
}

// GetSurroundings returns a range surrounding the position by amount in every
// direction.
func (pos ChunkPos) GetSurroundings(amount int) ChunkRange {
	minx := pos.X - amount
	maxx := pos.X + amount
	miny := pos.Y - amount
	maxy := pos.Y + amount
	mink := pos.Z - amount
	maxk := pos.Z + amount
	return ChunkRange{
		Min: ChunkPos{minx, miny, mink},
		Max: ChunkPos{maxx, maxy, maxk},
	}
}

// Mul returns this ChunkPos multiplied by a scalar s.
func (pos ChunkPos) Mul(s int) ChunkPos {
	return ChunkPos{
		X: pos.X * s,
		Y: pos.Y * s,
		Z: pos.Z * s,
	}
}

func (pos ChunkPos) AsVec3() glm.Vec3 {
	return glm.Vec3{
		float32(pos.X),
		float32(pos.Y),
		float32(pos.Z),
	}
}

func (c *Chunk) AsVoxelPos() VoxelPos {
	scaled := c.Pos.Mul(c.size)
	return VoxelPos{
		X: scaled.X,
		Y: scaled.Y,
		Z: scaled.Z,
	}
}

// ChunkRange is the range of chunks between Min and Max.
type ChunkRange struct {
	Min ChunkPos
	Max ChunkPos
}

// ForEach executes the given function on every position in the this ChunkRange.
func (rng ChunkRange) ForEach(fn func(pos ChunkPos)) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for y := rng.Min.Y; y <= rng.Max.Y; y++ {
			for z := rng.Min.Z; z <= rng.Max.Z; z++ {
				fn(ChunkPos{X: x, Y: y, Z: z})
			}
		}
	}
}

// Contains returns whether this ChunkRange contains the given pos.
func (rng ChunkRange) Contains(pos ChunkPos) bool {
	if pos.X < rng.Min.X || pos.X > rng.Max.X {
		return false
	}
	if pos.Y < rng.Min.Y || pos.Y > rng.Max.Y {
		return false
	}
	if pos.Z < rng.Min.Z || pos.Z > rng.Max.Z {
		return false
	}
	return true
}

// Chunk manages a size X height X size region of voxels.
type Chunk struct {
	Pos      ChunkPos
	flatData []float32
	objs     *voxgl.Object
	root     *Octree
	dirty    bool
	size     int
	modified bool
}

const VertSize = 4

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size int, pos ChunkPos, gen Generator) *Chunk {
	flatData := make([]float32, size*size*size*VertSize)
	chunk := &Chunk{
		Pos:      pos,
		flatData: flatData,
		size:     size,
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				x := chunk.AsVoxelPos().X + i
				y := chunk.AsVoxelPos().Y + j
				z := chunk.AsVoxelPos().Z + k
				vox := gen.GenerateAt(x, y, z)
				chunk.SetVoxel(vox)
			}
		}
	}
	return chunk
}

// NewChunkLoaded returns a pre-loaded chunk
func NewChunkLoaded(size int, pos ChunkPos, flatData []int32) *Chunk {
	chunk := &Chunk{
		Pos:      pos,
		flatData: make([]float32, 4*size*size*size),
		size:     size,
	}
	// rebuild octree
	maxIdx := 4 * size * size * size
	for i := 0; i < maxIdx; i += 4 {
		vbits := flatData[i+3]
		v := Voxel{
			Pos: VoxelPos{
				int(flatData[i]),
				int(flatData[i+1]),
				int(flatData[i+2]),
			},
			AdjMask: AdjacentMask(vbits & int32(AdjacentAll)),
			Btype:   BlockType(vbits & ^int32(AdjacentAll)),
		}
		if v.Btype != Air {
			chunk.root = chunk.root.AddLeaf(&v)
		}
		chunk.flatData[i] = float32(flatData[i])
		chunk.flatData[i+1] = float32(flatData[i+1])
		chunk.flatData[i+2] = float32(flatData[i+2])
		chunk.flatData[i+3] = float32(vbits)

	}
	chunk.modified = true
	chunk.dirty = true
	return chunk
}

func (c *Chunk) SetObjs(objs *voxgl.Object) {
	c.objs = objs
}

func (c *Chunk) GetRoot() *Octree {
	return c.root
}

func (c *Chunk) GetFlatData() []float32 {
	return c.flatData
}

// IsWithinChunk returns whether the position is within the chunk
func (c *Chunk) IsWithinChunk(pos VoxelPos) bool {
	if pos.X < c.AsVoxelPos().X || pos.Y < c.AsVoxelPos().Y || pos.Z < c.AsVoxelPos().Z {
		//pos is below x, y, or z chunk bounds
		return false
	}
	if pos.X >= c.AsVoxelPos().X+c.size || pos.Y >= c.AsVoxelPos().Y+c.size ||
		pos.Z >= c.AsVoxelPos().Z+c.size {
		// pos is above x, y, or z chunk bounds
		return false
	}
	return true
}

type BlockType int32

const (
	Air BlockType = iota
	Dirt
	Grass
	Labeled
)

// AdjacentMask indicates which in which directions there are adjacent voxels.
type AdjacentMask int32

const (
	AdjacentFront  AdjacentMask = 0b00000001 // The voxel has a backward adjacency.
	AdjacentBack   AdjacentMask = 0b00000010 // The voxel has a forward adjacency.
	AdjacentBottom AdjacentMask = 0b00000100 // The voxel has a bottom adjacency.
	AdjacentTop    AdjacentMask = 0b00001000 // The voxel has a top adjacency.
	AdjacentLeft   AdjacentMask = 0b00010000 // The voxel has a right adjacency.
	AdjacentRight  AdjacentMask = 0b00100000 // The voxel has a left adjacency.

	AdjacentX   = AdjacentRight | AdjacentLeft      // The voxel has adjacencies in the +/-x directions.
	AdjacentY   = AdjacentTop | AdjacentBottom      // The voxel has adjacencies in the +/-y directions.
	AdjacentZ   = AdjacentBack | AdjacentFront      // The voxel has adjacencies in the +/-z directions.
	AdjacentAll = AdjacentX | AdjacentY | AdjacentZ // The voxel has adjacencies in all directions.
)

func (c *Chunk) SetModified() {
	c.modified = true
}

// SetVoxel updates a voxel's variables in the chunk, if it exists
func (c *Chunk) SetVoxel(v *Voxel) {
	if !c.IsWithinChunk(v.Pos) {
		log.Debugf("%v is not within %v", v, c.AsVoxelPos())
		return
	}
	x, y, z := float32(v.Pos.X), float32(v.Pos.Y), float32(v.Pos.Z)
	localPos := v.Pos.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	vbits := float32(int32(v.AdjMask) | int32(v.Btype<<6))
	off := (i + j*c.size*c.size + k*c.size) * VertSize
	if off%VertSize != 0 {
		panic("offset not divisible by VertSize")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	c.flatData[off] = x
	c.flatData[off+1] = y
	c.flatData[off+2] = z
	c.flatData[off+3] = vbits

	if v.Btype != Air { // TODO return at top of function?
		c.root = c.root.AddLeaf(v)
	}
	c.dirty = true
}

// AddAdjacency adds adjacency to a voxel
func (c *Chunk) AddAdjacency(v VoxelPos, adjMask AdjacentMask) {
	if !c.IsWithinChunk(v) {
		log.Debugf("%v is not within %v", v, c.AsVoxelPos())
		return
	}
	localPos := v.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	off := (i + j*c.size*c.size + k*c.size) * VertSize
	vbits := int32(c.flatData[off+3]) | int32(adjMask)
	if off%VertSize != 0 {
		panic("offset not divisible by VertSize")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	c.flatData[off+3] = float32(vbits)
	c.dirty = true
	c.modified = true
}

// RemoveAdjacency remove adjacency from a voxel
func (c *Chunk) RemoveAdjacency(v VoxelPos, adjMask AdjacentMask) {
	if !c.IsWithinChunk(v) {
		log.Debugf("%v is not within %v", v, c.AsVoxelPos())
		return
	}
	localPos := v.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	off := (i + j*c.size*c.size + k*c.size) * VertSize
	vbits := int32(c.flatData[off+3]) & ^int32(adjMask)
	if off%VertSize != 0 {
		panic("offset not divisible by VertSize")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	c.flatData[off+3] = float32(vbits)
	c.dirty = true
	c.modified = true
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
