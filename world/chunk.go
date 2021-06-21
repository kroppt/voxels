package world

import (
	"fmt"

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
	Pos       ChunkPos
	flatData  []float32
	lights    map[VoxelPos]struct{}
	objs      *voxgl.Object
	root      *Octree
	dirty     bool
	size      int
	modified  bool
	empty     bool
	lightPass int32
}

const VertSize = 5
const CacheVertSize = 4

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size int, pos ChunkPos, gen Generator) *Chunk {
	flatData := make([]float32, size*size*size*VertSize)
	chunk := &Chunk{
		Pos:      pos,
		flatData: flatData,
		size:     size,
		empty:    true,
		lights:   make(map[VoxelPos]struct{}),
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
	chunk.lightingAlgo()
	return chunk
}

// NewChunkLoaded returns a pre-loaded chunk
func NewChunkLoaded(size int, pos ChunkPos, flatData []int32) *Chunk {
	chunk := &Chunk{
		Pos:      pos,
		flatData: make([]float32, VertSize*size*size*size),
		size:     size,
		empty:    true,
		lights:   make(map[VoxelPos]struct{}),
	}
	maxIdx := CacheVertSize * size * size * size
	for i := 0; i < maxIdx; i += CacheVertSize {
		vbits := flatData[i+3]
		adjMask, btype := SeparateVbits(vbits)
		v := Voxel{
			Pos: VoxelPos{
				int(flatData[i]),
				int(flatData[i+1]),
				int(flatData[i+2]),
			},
			AdjMask: adjMask,
			Btype:   btype,
		}
		chunk.SetVoxel(&v)
	}
	chunk.lightingAlgo()
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
	Corrupted
	Stone
	Light
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

	AdjacentX    = AdjacentRight | AdjacentLeft      // The voxel has adjacencies in the +/-x directions.
	AdjacentY    = AdjacentTop | AdjacentBottom      // The voxel has adjacencies in the +/-y directions.
	AdjacentZ    = AdjacentBack | AdjacentFront      // The voxel has adjacencies in the +/-z directions.
	AdjacentAll  = AdjacentX | AdjacentY | AdjacentZ // The voxel has adjacencies in all directions.
	AdjacentNone = 0
)

type LightMask uint32

const (
	MaxLightValue        = 5
	ExploredMask  uint32 = 1 << 30

	LightFront  LightMask = 0b1111           // The voxel's front face lighting bits.
	LightBack   LightMask = LightFront << 4  // The voxel's back face lighting bits.
	LightBottom LightMask = LightFront << 8  // The voxel's bottom face lighting bits.
	LightTop    LightMask = LightFront << 12 // The voxel's top face lighting bits.
	LightLeft   LightMask = LightFront << 16 // The voxel's left face lighting bits.
	LightRight  LightMask = LightFront << 20 // The voxel's right face lighting bits.
	LightValue  LightMask = LightFront << 24
	LightAll              = LightFront | LightBack | LightBottom | LightTop | LightLeft | LightRight
)

func (c *Chunk) SetModified() {
	c.modified = true
}

func (c *Chunk) lightingAlgo() {
	for lightPos := range c.lights {
		// TODO extract hard coded max lighting value, especially when
		// we might want to change this maximum to increase the range
		c.lightFrom(lightPos, MaxLightValue)
		// TODO avoid resetting by storing identifier for light block to know if it owns the explored bit
		// identifier could be local coordinate
		for i := 0; i < c.size; i++ {
			for j := 0; j < c.size; j++ {
				for k := 0; k < c.size; k++ {
					p := VoxelPos{
						X: c.AsVoxelPos().X + i,
						Y: c.AsVoxelPos().Y + j,
						Z: c.AsVoxelPos().Z + k,
					}
					v := c.GetVoxel(p)
					v.SetLightValue(0, LightValue)
					c.SetVoxel(&v)
				}
			}
		}
	}
}

func (c *Chunk) lightFrom(p VoxelPos, value uint32) {
	if value < 0 || value > MaxLightValue || !c.IsWithinChunk(p) {
		panic("improper usage: bad lighting value or p outside chunk")
	}
	if value == 0 {
		return
	}
	currBlock := c.GetVoxel(p)
	if currBlock.GetLightValue(LightLeft) >= value {
		return
	}
	if currBlock.Btype == Air {
		currBlock.SetLightValue(value, LightLeft)
		c.SetVoxel(&currBlock)
	}
	dirs := []struct {
		off  VoxelPos
		face LightMask
	}{
		{VoxelPos{-1, 0, 0}, LightRight},
		{VoxelPos{1, 0, 0}, LightLeft},
		{VoxelPos{0, -1, 0}, LightTop},
		{VoxelPos{0, 1, 0}, LightBottom},
		{VoxelPos{0, 0, -1}, LightBack},
		{VoxelPos{0, 0, 1}, LightFront},
	}
	for _, dir := range dirs {
		offP := p.Add(dir.off)
		if c.IsWithinChunk(offP) {
			adjBlock := c.GetVoxel(offP)
			if adjBlock.Btype == Air {
				// continue search down open path
				c.lightFrom(offP, value-1)
			} else {
				// path is blocked off, apply lighting value to the face
				if adjBlock.GetLightValue(dir.face) < value {
					log.Debugf("Setting face %v of %v to value %v", dir.face, offP, value)
					adjBlock.SetLightValue(value, dir.face)
					c.SetVoxel(&adjBlock)
				}
			}
		} else {
			// TODO lighting across chunks
			// just within for now
		}
	}
}

func (c *Chunk) GetVoxel(pos VoxelPos) Voxel {
	if !c.IsWithinChunk(pos) {
		panic(fmt.Sprintf("%v is not within %v", pos, c.AsVoxelPos()))
	}
	localPos := pos.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	off := (i + j*c.size*c.size + k*c.size) * VertSize
	vbits := int32(c.flatData[off+3])
	adjMask, btype := SeparateVbits(vbits)
	lbits := uint32(c.flatData[off+4])
	explored, lightBits := SeparateLbits(lbits)
	return Voxel{
		Pos: VoxelPos{
			X: int(c.flatData[off]),
			Y: int(c.flatData[off+1]),
			Z: int(c.flatData[off+2]),
		},
		AdjMask:   adjMask,
		Btype:     btype,
		Explored:  explored,
		LightBits: lightBits,
	}
}

// SetVoxel updates a voxel's variables in the chunk, if it exists
func (c *Chunk) SetVoxel(v *Voxel) {
	if !c.IsWithinChunk(v.Pos) {
		panic(fmt.Sprintf("%v is not within %v", v, c.AsVoxelPos()))
	}
	x, y, z := float32(v.Pos.X), float32(v.Pos.Y), float32(v.Pos.Z)
	localPos := v.Pos.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	off := (i + j*c.size*c.size + k*c.size) * VertSize
	if off%VertSize != 0 {
		panic("offset not divisible by VertSize")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	oldVox := c.GetVoxel(v.Pos)
	if v.Btype == Light {
		c.lights[v.Pos] = struct{}{}
	} else if oldVox.Btype == Light && v.Btype != Light {
		delete(c.lights, v.Pos)
	}

	c.flatData[off] = x
	c.flatData[off+1] = y
	c.flatData[off+2] = z
	c.flatData[off+3] = float32(v.GetVbits())
	c.flatData[off+4] = float32(v.GetLbits())

	if v.Btype != Air {
		c.root = c.root.AddLeaf(v)
		c.empty = false
	}

	c.dirty = true
}

// AddAdjacency adds adjacency to a voxel
func (c *Chunk) AddAdjacency(v VoxelPos, adjMask AdjacentMask) {
	if !c.IsWithinChunk(v) {
		panic(fmt.Sprintf("%v is not within %v", v, c.AsVoxelPos()))
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
		panic(fmt.Sprintf("%v is not within %v", v, c.AsVoxelPos()))
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
func (c *Chunk) Render(cam *Camera) {
	if c.dirty {
		c.objs.SetData(c.flatData)
		c.dirty = false
	}
	if c.empty {
		return
	}
	if cam.IsWithinFrustum(c.AsVoxelPos().AsVec3(), float32(c.size), float32(c.size), float32(c.size)) {
		c.objs.Render()
	}
}

func (c *Chunk) Destroy() {
	c.objs.Destroy()
}
