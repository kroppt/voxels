package world

import (
	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/voxgl"
)

// Chunk manages a size X size X height region of voxels.
type Chunk struct {
	Pos      glm.Vec2
	flatData []float32
	voxels   [][][]*Voxel
	objs     *voxgl.Object
	root     *Octree
	dirty    bool
}

// NewChunk returns a new Chunk shaped as size X size X height
// TODO height should be in y?
func NewChunk(size, height int32, pos glm.Vec2) *Chunk {
	voxels := make([][][]*Voxel, height)
	for i := range voxels {
		voxels[i] = make([][]*Voxel, size)
		for j := range voxels[i] {
			voxels[i][j] = make([]*Voxel, size)
		}
	}
	objs, err := voxgl.NewColoredObject(nil)
	if err != nil {
		panic("failed to make NewColoredObject for chunk")
	}
	return &Chunk{
		Pos:    pos,
		objs:   objs,
		voxels: voxels,
	}
}

func (c *Chunk) GetRoot() *Octree {
	return c.root
}

// AddVoxel adds a voxel to the chunk, updating all data structures
func (c *Chunk) AddVoxel(v *Voxel) {
	i, j, k := v.Pos[0], v.Pos[1], v.Pos[2]
	r, g, b, a := v.Col[0], v.Col[1], v.Col[2], v.Col[3]
	c.flatData = append(c.flatData, i, j, k, r, g, b, a)
	// TODO translate to relative chunk coordinate
	// c.voxels[int32(i)][int32(j)][int32(k)] = v
	c.root = c.root.AddLeaf(v)
	c.dirty = true
}

func (c *Chunk) Render() {
	if c.dirty {
		c.objs.SetData(c.flatData)
		c.dirty = false
	}
	c.objs.Render()
}
