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

func New(chPos Position, chSize uint32) Chunk {
	if chSize == 0 {
		panic("chunk size cannot be 0")
	}

	flatData := make([]float32, vertSize*chSize*chSize*chSize)
	size := int32(chSize)
	chunk := Chunk{
		pos:      chPos,
		size:     chSize,
		flatData: flatData,
	}
	for x := chPos.X * size; x < chPos.X*size+size; x++ {
		for y := chPos.Y * size; y < chPos.Y*size+size; y++ {
			for z := chPos.Z * size; z < chPos.Z*size+size; z++ {
				off := chunk.voxelPosToDataOffset(VoxelCoordinate{x, y, z})
				chunk.flatData[off] = float32(x)
				chunk.flatData[off+1] = float32(y)
				chunk.flatData[off+2] = float32(z)
			}
		}
	}
	return chunk
}

func (c Chunk) Position() Position {
	return c.pos
}

func (c Chunk) Size() uint32 {
	return c.size
}

func (c Chunk) GetFlatData() []float32 {
	return c.flatData
}

func (c Chunk) isOutOfBounds(vpos VoxelCoordinate) bool {
	return VoxelCoordToChunkCoord(vpos, c.size) != c.pos
}

func (c Chunk) voxelPosToDataOffset(vpos VoxelCoordinate) int32 {
	if c.isOutOfBounds(vpos) {
		panic("voxel position is out of chunk bounds")
	}
	size := int32(c.size)
	i := vpos.X - c.pos.X*size
	j := vpos.Y - c.pos.Y*size
	k := vpos.Z - c.pos.Z*size
	return (i + j*size*size + k*size) * vertSize
}

func (c Chunk) SetBlockType(vpos VoxelCoordinate, btype BlockType) {
	off := c.voxelPosToDataOffset(vpos)
	vbits := uint32(c.flatData[off+3])
	btypeBits := uint32(btype) << 6
	c.flatData[off+3] = float32(vbits&(^btypeMask) | btypeBits)
}

func (c Chunk) BlockType(vpos VoxelCoordinate) BlockType {
	off := c.voxelPosToDataOffset(vpos)
	vbits := uint32(c.flatData[off+3])
	return BlockType(vbits >> 6)
}

func (c Chunk) SetAdjacency(vpos VoxelCoordinate, adj AdjacentMask) {
	if adj > AdjacentAll {
		panic("invalid adj mask")
	}
	off := c.voxelPosToDataOffset(vpos)
	vbits := uint32(c.flatData[off+3])
	adjBits := uint32(adj)
	c.flatData[off+3] = float32(vbits&(^uint32(AdjacentAll)) | adjBits)
}

func (c Chunk) Adjacency(vpos VoxelCoordinate) AdjacentMask {
	off := c.voxelPosToDataOffset(vpos)
	vbits := uint32(c.flatData[off+3])
	return AdjacentMask(vbits & uint32(AdjacentAll))
}

func (c Chunk) SetLighting(vpos VoxelCoordinate, face LightFace, intensity uint32) {
	if intensity > 15 {
		panic("light intensity too high")
	}
	if face > bitsPerMask*5 {
		panic("invalid light face specified")
	}
	off := c.voxelPosToDataOffset(vpos)
	lbits := uint32(c.flatData[off+4])
	c.flatData[off+4] = float32(lbits&(^uint32(lightAll)) | (intensity << uint32(face)))
}

func (c Chunk) Lighting(vpos VoxelCoordinate, face LightFace) uint32 {
	if face > bitsPerMask*5 {
		panic("invalid light face specified")
	}
	off := c.voxelPosToDataOffset(vpos)
	mask := uint32(0b1111 << int(face))
	lbits := uint32(c.flatData[off+4])
	return (lbits & mask) >> face
}

func VoxelCoordToChunkCoord(pos VoxelCoordinate, chunkSize uint32) Position {
	if chunkSize == 0 {
		panic("chunk size 0 is invalid")
	}
	x, y, z := pos.X, pos.Y, pos.Z
	size := int32(chunkSize)
	if pos.X < 0 {
		x++
	}
	if pos.Y < 0 {
		y++
	}
	if pos.Z < 0 {
		z++
	}
	x /= size
	y /= size
	z /= size
	if pos.X < 0 {
		x--
	}
	if pos.Y < 0 {
		y--
	}
	if pos.Z < 0 {
		z--
	}
	return Position{X: x, Y: y, Z: z}
}
