package chunk

type Chunk struct {
	pos      ChunkCoordinate
	size     uint32
	flatData []float32
}

type ChunkCoordinate struct {
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
	BlockTypeGrass
	BlockTypeGrassSides
	BlockTypeLabeled
	BlockTypeCorrupted
	BlockTypeStone
	BlockTypeLight
)
const LargestVbits = uint32(BlockTypeDirt)<<6 | uint32(AdjacentAll)

const VertSize = 5
const BytesPerElement = 4
const adjacencyMask = 0x0000003F
const btypeMask uint32 = 0xFFFFFFC0

// AdjacentMask indicates which in which directions there are adjacent voxels.
type AdjacentMask uint32

const (
	AdjacentFront  AdjacentMask = 0b00000001 // There is a voxel in front of it (-Z).
	AdjacentBack   AdjacentMask = 0b00000010 // There is a voxel behind it (+Z).
	AdjacentBottom AdjacentMask = 0b00000100 // There is a voxel below it (-Y).
	AdjacentTop    AdjacentMask = 0b00001000 // There is a voxel above it (+Y).
	AdjacentLeft   AdjacentMask = 0b00010000 // There is a voxel left of it (-X).
	AdjacentRight  AdjacentMask = 0b00100000 // There is a voxel right of it (+X).

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
	LightAll    uint32    = 0x00FFFFFF
)

func NewChunkEmpty(chPos ChunkCoordinate, chSize uint32) Chunk {
	if chSize == 0 {
		panic("chunk size cannot be 0")
	}

	flatData := make([]float32, VertSize*chSize*chSize*chSize)
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

func NewChunkFromData(data []float32, chSize uint32, chPos ChunkCoordinate) Chunk {
	if len(data) != int(VertSize*chSize*chSize*chSize) {
		panic("new chunk data has wrong size")
	}
	size := int32(chSize)
	ch := Chunk{
		pos:  chPos,
		size: chSize,
	}
	for x := chPos.X * size; x < chPos.X*size+size; x++ {
		for y := chPos.Y * size; y < chPos.Y*size+size; y++ {
			for z := chPos.Z * size; z < chPos.Z*size+size; z++ {
				off := ch.voxelPosToDataOffset(VoxelCoordinate{x, y, z})
				if data[off] != float32(x) {
					panic("invalid X coordinate in chunk data")
				}
				if data[off+1] != float32(y) {
					panic("invalid Y coordinate in chunk data")
				}
				if data[off+2] != float32(z) {
					panic("invalid Z coordinate in chunk data")
				}
				if data[off+3] > float32(LargestVbits) {
					panic("invalid vbits in chunk data")
				}
				if data[off+4] > float32(LightAll) {
					panic("invalid lighting bits in chunk data")
				}
			}
		}
	}
	ch.flatData = data
	return ch
}

func (c Chunk) ForEachVoxel(f func(VoxelCoordinate)) {
	size := int32(c.size)
	for x := c.pos.X * size; x < c.pos.X*size+size; x++ {
		for y := c.pos.Y * size; y < c.pos.Y*size+size; y++ {
			for z := c.pos.Z * size; z < c.pos.Z*size+size; z++ {
				f(VoxelCoordinate{x, y, z})
			}
		}
	}
}

func (c Chunk) Position() ChunkCoordinate {
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
	return (i + j*size + k*size*size) * VertSize
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
	c.flatData[off+4] = float32(lbits&(^uint32(LightAll)) | (intensity << uint32(face)))
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

func (c Chunk) Vbits(vpos VoxelCoordinate) uint32 {
	return uint32(c.BlockType(vpos))<<6 | uint32(c.Adjacency(vpos))
}

func VoxelCoordToChunkCoord(pos VoxelCoordinate, chunkSize uint32) ChunkCoordinate {
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
	return ChunkCoordinate{X: x, Y: y, Z: z}
}
