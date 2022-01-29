package chunk

type Chunk struct {
	pos      Position
	size     uint32
	flatData []float32
}

type Position struct {
	X int32
	Y int32
	Z int32
}

type VoxelCoordinate struct {
	X int32
	Y int32
	Z int32
}

type BlockType uint32

const (
	BlockTypeAir BlockType = iota
	BlockTypeDirt
)

const vertSize = 5
const adjacencyMask = 0x0000003F
const btypeMask uint32 = 0xFFFFFFC0

// AdjacentMask indicates which in which directions there are adjacent voxels.
type AdjacentMask uint32

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
	AdjacentNone = 0                                 // The voxel has no adjacencies.
)

type LightFace uint32

const (
	maxLightValue = 8
	bitsPerMask   = 4

	LightFront  LightFace = 0               // The voxel's front face lighting bits.
	LightBack             = bitsPerMask     // The voxel's back face lighting bits.
	LightBottom           = bitsPerMask * 2 // The voxel's bottom face lighting bits.
	LightTop              = bitsPerMask * 3 // The voxel's top face lighting bits.
	LightLeft             = bitsPerMask * 4 // The voxel's left face lighting bits.
	LightRight            = bitsPerMask * 5 // The voxel's right face lighting bits.
	lightAll              = 0x00FFFFFF
)

func New(pos Position, size uint32) Chunk {
	if size == 0 {
		panic("chunk size cannot be 0")
	}
	return Chunk{
		pos:      pos,
		size:     size,
		flatData: make([]float32, vertSize*size*size*size),
	}
}

func (c Chunk) Position() Position {
	return c.pos
}

func (c Chunk) Size() uint32 {
	return c.size
}

func (c Chunk) isOutOfBounds(vpos VoxelCoordinate) bool {
	chMinX := c.pos.X * int32(c.size)
	chMaxX := chMinX + int32(c.size)
	chMinY := c.pos.Y * int32(c.size)
	chMaxY := chMinY + int32(c.size)
	chMinZ := c.pos.Z * int32(c.size)
	chMaxZ := chMinZ + int32(c.size)
	if vpos.X < chMinX || vpos.Y < chMinY || vpos.Z < chMinZ {
		return true
	}
	if vpos.X >= chMaxX || vpos.Y >= chMaxY || vpos.Z >= chMaxZ {
		return true
	}
	return false
}

func (c Chunk) SetBlockType(vpos VoxelCoordinate, btype BlockType) {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	vbits := uint32(c.flatData[off+3])
	btypeBits := uint32(btype) << 6
	c.flatData[off+3] = float32(vbits&(^btypeMask) | btypeBits)
}

func (c Chunk) BlockType(vpos VoxelCoordinate) BlockType {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	vbits := uint32(c.flatData[off+3])
	return BlockType(vbits >> 6)
}

func (c Chunk) SetAdjacency(vpos VoxelCoordinate, adj AdjacentMask) {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	if adj > AdjacentAll {
		panic("invalid adj mask")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	vbits := uint32(c.flatData[off+3])
	adjBits := uint32(adj)
	c.flatData[off+3] = float32(vbits&(^uint32(AdjacentAll)) | adjBits)
}

func (c Chunk) Adjacency(vpos VoxelCoordinate) AdjacentMask {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	vbits := uint32(c.flatData[off+3])
	return AdjacentMask(vbits & uint32(AdjacentAll))
}

func (c Chunk) SetLighting(vpos VoxelCoordinate, face LightFace, intensity uint32) {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	if intensity > 15 {
		panic("light intensity too high")
	}
	if face > bitsPerMask*5 {
		panic("invalid light face specified")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	lbits := uint32(c.flatData[off+4])
	c.flatData[off+4] = float32(lbits&(^uint32(lightAll)) | (intensity << uint32(face)))
}

func (c Chunk) Lighting(vpos VoxelCoordinate, face LightFace) uint32 {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	if face > bitsPerMask*5 {
		panic("invalid light face specified")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	off := (i + j*size*size + k*size) * vertSize
	mask := uint32(0b1111 << int(face))
	lbits := uint32(c.flatData[off+4])
	return (lbits & mask) >> face
}
