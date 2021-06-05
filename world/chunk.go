package world

import (
	"math/rand"
	"time"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/voxgl"
)

// Chunk manages a size X height X size region of voxels.
type Chunk struct {
	Pos      glm.Vec2
	flatData []float32
	voxels   [][][]*Voxel
	objs     *voxgl.Object
	root     *Octree
	dirty    bool
}

// NewChunk returns a new Chunk shaped as size X height X size.
func NewChunk(size, height int, chunkPos glm.Vec2) *Chunk {
	voxels := make([][][]*Voxel, size)
	for i := range voxels {
		voxels[i] = make([][]*Voxel, height)
		for j := range voxels[i] {
			voxels[i][j] = make([]*Voxel, size)
		}
	}
	objs, err := voxgl.NewColoredObject(nil)
	if err != nil {
		panic("failed to make NewColoredObject for chunk")
	}
	chunk := &Chunk{
		Pos:    chunkPos.Mul(float32(size)),
		objs:   objs,
		voxels: voxels,
	}
	rand.Seed(time.Now().UnixNano())
	for i := range voxels {
		for j := range voxels[i] {
			for k := range voxels[i][j] {
				x := chunk.Pos.X() + float32(i)
				y := float32(j)
				z := chunk.Pos.Y() + float32(k)
				r, g, b := rand.Float32(), rand.Float32(), rand.Float32()
				v := Voxel{
					Pos: glm.Vec3{x, y, z},
					Col: glm.Vec4{r, g, b, 1.0},
				}
				chunk.AddVoxel(&v)
			}
		}
	}
	return chunk
}

func (c *Chunk) GetRoot() *Octree {
	return c.root
}

func (c *Chunk) getRelativeIndices(pos glm.Vec3) (i int, j int, k int) {
	return int(pos.X() - c.Pos.X()), int(pos.Y()), int(pos.Z() - c.Pos.Y())
}

// AddVoxel adds a voxel to the chunk, updating all data structures.
func (c *Chunk) AddVoxel(v *Voxel) {
	i, j, k := v.Pos.X(), v.Pos.Y(), v.Pos.Z()
	r, g, b, a := v.Col[0], v.Col[1], v.Col[2], v.Col[3]
	c.flatData = append(c.flatData, i, j, k, r, g, b, a)
	x, y, z := c.getRelativeIndices(v.Pos)
	c.voxels[x][y][z] = v
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
