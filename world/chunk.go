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
	X int
	Z int
}

// GetSurroundings returns a range surrounding the position by amount in every
// direction.
func (pos ChunkPos) GetSurroundings(amount int) ChunkRange {
	minx := pos.X - amount
	maxx := pos.X + amount
	mink := pos.Z - amount
	maxk := pos.Z + amount
	return ChunkRange{
		Min: ChunkPos{minx, mink},
		Max: ChunkPos{maxx, maxk},
	}
}

// Mul returns this ChunkPos multiplied by a scalar s.
func (pos ChunkPos) Mul(s int) ChunkPos {
	return ChunkPos{
		X: pos.X * s,
		Z: pos.Z * s,
	}
}

func (pos ChunkPos) AsVec3() glm.Vec3 {
	return glm.Vec3{
		float32(pos.X),
		// TODO un hack this, chunks just so happen to start at y=0
		float32(0.0),
		float32(pos.Z),
	}
}

func (c *Chunk) AsVoxelPos() VoxelPos {
	scaled := c.Pos.Mul(c.size)
	return VoxelPos{
		X: scaled.X,
		Y: 0, // TODO hacked
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
		for z := rng.Min.Z; z <= rng.Max.Z; z++ {
			fn(ChunkPos{X: x, Z: z})
		}
	}
}

// Contains returns whether this ChunkRange contains the given pos.
func (rng ChunkRange) Contains(pos ChunkPos) bool {
	if pos.X < rng.Min.X || pos.X > rng.Max.X {
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
	height   int
}

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size, height int, pos ChunkPos, chunkChan chan *Chunk) {
	vertSize := 8
	flatData := make([]float32, size*size*height*vertSize)
	// layout 4+4=8 hard coded in here too
	chunk := &Chunk{
		Pos:      pos,
		flatData: flatData,
		size:     size,
		height:   height,
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < size; k++ {
				x := chunk.AsVoxelPos().X + i
				y := j
				z := chunk.AsVoxelPos().Z + k
				r, g, b := rand.Float32(), rand.Float32(), rand.Float32()
				v := Voxel{
					Pos:   VoxelPos{x, y, z},
					Color: Color{r, g, b, 1.0},
				}
				// left right top bottom forward backward
				// a worldgenerator should be asked for this info
				leftmask := 0x20
				rightmask := 0x10
				topmask := 0x08
				bottommask := 0x04
				forwardmask := 0x02
				backwardmask := 0x01
				code := 0
				if i == 0 {
					code += leftmask
				}
				if i == size-1 {
					code += rightmask
				}
				if j == 0 {
					code += bottommask
				}
				if j == height-1 {
					code += topmask
				}
				if k == 0 {
					code += forwardmask
				}
				if k == size-1 {
					code += backwardmask
				}
				chunk.SetVoxel(&v, float32(code))
			}
		}
	}
	chunkChan <- chunk
}

func (c *Chunk) SetObjs(objs *voxgl.Object) {
	c.objs = objs
}

func (c *Chunk) GetRoot() *Octree {
	return c.root
}

// IsWithinChunk returns whether the position is within the chunk
func (c *Chunk) IsWithinChunk(pos VoxelPos) bool {
	if pos.X < c.AsVoxelPos().X || pos.Z < c.AsVoxelPos().Z {
		//pos is below x or z chunk bounds
		return false
	}
	if pos.X >= c.AsVoxelPos().X+c.size || pos.Z >= c.AsVoxelPos().Z+c.size {
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
func (c *Chunk) SetVoxel(v *Voxel, adjaCode float32) {
	if !c.IsWithinChunk(v.Pos) {
		log.Debugf("%v is not within %v", v, c.AsVoxelPos())
		return
	}
	x, y, z := float32(v.Pos.X), float32(v.Pos.Y), float32(v.Pos.Z)
	localPos := v.Pos.AsLocalChunkPos(*c)
	i, j, k := localPos.X, localPos.Y, localPos.Z
	r, g, b, a := v.Color.R, v.Color.G, v.Color.B, v.Color.A
	off := (i + j*c.size*c.size + k*c.size) * 8
	if off%8 != 0 {
		panic("offset not divisible by 8")
	}
	if off >= len(c.flatData) || off < 0 {
		panic("offset out of bounds")
	}

	c.flatData[off] = x
	c.flatData[off+1] = y
	c.flatData[off+2] = z
	c.flatData[off+3] = adjaCode
	c.flatData[off+4] = r
	c.flatData[off+5] = g
	c.flatData[off+6] = b
	c.flatData[off+7] = a

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
