package world

import (
	"math/rand"
	"time"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/voxgl"
)

// Chunk manages a size X height X size region of voxels.
type Chunk struct {
	Pos      glm.Vec2
	flatData []float32
	objs     *voxgl.Object
	root     *Octree
	dirty    bool
	size     int
	height   int
}

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size, height int, chunkPos glm.Vec2) *Chunk {
	vertSize := 7
	flatData := make([]float32, size*size*height*vertSize)
	// layout 3+4=7 hard coded in here too
	objs, err := voxgl.NewColoredObject(nil)
	if err != nil {
		panic("failed to make NewColoredObject for chunk")
	}
	chunk := &Chunk{
		Pos:      chunkPos.Mul(float32(size)),
		objs:     objs,
		flatData: flatData,
		size:     size,
		height:   height,
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < size; k++ {
				x := chunk.Pos.X() + float32(i)
				y := float32(j)
				z := chunk.Pos.Y() + float32(k)
				r, g, b := rand.Float32(), rand.Float32(), rand.Float32()
				v := Voxel{
					Pos: glm.Vec3{x, y, z},
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
func (c *Chunk) GetRelativeIndices(pos glm.Vec3) (int, int, int) {
	return int(pos.X() - c.Pos.X()), int(pos.Y()), int(pos.Z() - c.Pos.Y())
}

// IsWithinChunk returns whether the position is within the chunk
func (c *Chunk) IsWithinChunk(pos glm.Vec3) bool {
	if pos.X() < c.Pos.X() || pos.Z() < c.Pos.Y() {
		//pos is below x or z chunk bounds
		return false
	}
	if pos.X() >= c.Pos.X()+float32(c.size) || pos.Z() >= c.Pos.Y()+float32(c.size) {
		// pos is above x or z chunk bounds
		return false
	}
	if pos.Y() < 0 || pos.Y() >= float32(c.height) {
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
	x, y, z := v.Pos.X(), v.Pos.Y(), v.Pos.Z()
	i, j, k := c.GetRelativeIndices(v.Pos)
	r, g, b, a := v.Col[0], v.Col[1], v.Col[2], v.Col[3]
	off := (i + j*c.size*c.height + k*c.size) * 7
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
