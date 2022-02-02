package world

import (
	"math"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/repositories/settings"
)

type Generator interface {
	GenerateChunk(chunk.ChunkCoordinate) chunk.Chunk
}

type FnGenerator struct {
	FnGenerateChunk func(chunk.ChunkCoordinate) chunk.Chunk
}

func (fn *FnGenerator) GenerateChunk(pos chunk.ChunkCoordinate) chunk.Chunk {
	if fn.FnGenerateChunk != nil {
		return fn.FnGenerateChunk(pos)
	}
	return chunk.NewChunkEmpty(pos, 1)
}

type TestGenerator struct {
	settingsRepo settings.Interface
}

func NewTestGenerator(settingsRepo settings.Interface) *TestGenerator {
	if settingsRepo == nil {
		panic("test world generator missing settings repo")
	}
	return &TestGenerator{
		settingsRepo: settingsRepo,
	}
}

func (gen *TestGenerator) GenerateChunk(key chunk.ChunkCoordinate) chunk.Chunk {
	newChunk := chunk.NewChunkEmpty(key, gen.settingsRepo.GetChunkSize())
	if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
		newChunk.SetBlockType(chunk.VoxelCoordinate{
			X: 0,
			Y: 0,
			Z: 0,
		}, chunk.BlockTypeLabeled)
	}
	return newChunk
}

type FlatWorldGenerator struct {
	settingsRepo settings.Interface
}

func NewFlatWorldGenerator(settingsRepo settings.Interface) *FlatWorldGenerator {
	if settingsRepo == nil {
		panic("flat world generator missing settings repo")
	}
	return &FlatWorldGenerator{
		settingsRepo: settingsRepo,
	}
}

type genVoxel struct {
	adjMask chunk.AdjacentMask
	bType   chunk.BlockType
}

func (gen *FlatWorldGenerator) GenerateChunk(chPos chunk.ChunkCoordinate) chunk.Chunk {
	size := int32(gen.settingsRepo.GetChunkSize())
	ch := chunk.NewChunkEmpty(chPos, uint32(size))
	for x := chPos.X * size; x < chPos.X*size+size; x++ {
		for y := chPos.Y * size; y < chPos.Y*size+size; y++ {
			for z := chPos.Z * size; z < chPos.Z*size+size; z++ {
				voxInfo := gen.generateAt(x, y, z)
				ch.SetAdjacency(chunk.VoxelCoordinate{X: x, Y: y, Z: z}, voxInfo.adjMask)
				ch.SetBlockType(chunk.VoxelCoordinate{X: x, Y: y, Z: z}, voxInfo.bType)
			}
		}
	}
	return ch
}

func (gen *FlatWorldGenerator) generateAt(x, y, z int32) *genVoxel {
	if y < 0 || y > 6 {
		return &genVoxel{
			adjMask: chunk.AdjacentNone,
			bType:   chunk.BlockTypeAir,
		}
	}
	if y == 0 {
		return &genVoxel{
			adjMask: chunk.AdjacentAll & ^chunk.AdjacentBottom,
			bType:   chunk.BlockTypeLabeled,
		}
	} else if y == 6 {
		if x == 3 && z == 3 {
			return &genVoxel{
				adjMask: chunk.AdjacentAll & ^chunk.AdjacentTop,
				bType:   chunk.BlockTypeLight,
			}
		} else {
			return &genVoxel{
				adjMask: chunk.AdjacentAll & ^chunk.AdjacentTop,
				bType:   chunk.BlockTypeGrass,
			}
		}
	} else if y == 1 || y == 2 {
		return &genVoxel{
			adjMask: chunk.AdjacentAll,
			bType:   chunk.BlockTypeCorrupted,
		}
	} else if y == 3 || y == 4 {
		return &genVoxel{
			adjMask: chunk.AdjacentAll,
			bType:   chunk.BlockTypeStone,
		}
	} else {
		return &genVoxel{
			adjMask: chunk.AdjacentAll,
			bType:   chunk.BlockTypeDirt,
		}
	}
}

type AlexWorldGenerator struct {
	settingsRepo settings.Interface
}

func NewAlexWorldGenerator(settingsRepo settings.Interface) *AlexWorldGenerator {
	if settingsRepo == nil {
		panic("alex world generator missing settings repo")
	}
	return &AlexWorldGenerator{
		settingsRepo: settingsRepo,
	}
}

func (gen *AlexWorldGenerator) GenerateChunk(chPos chunk.ChunkCoordinate) chunk.Chunk {
	size := int32(gen.settingsRepo.GetChunkSize())
	ch := chunk.NewChunkEmpty(chPos, uint32(size))
	for x := chPos.X * size; x < chPos.X*size+size; x++ {
		for y := chPos.Y * size; y < chPos.Y*size+size; y++ {
			for z := chPos.Z * size; z < chPos.Z*size+size; z++ {
				voxInfo := gen.generateAt(x, y, z)
				ch.SetAdjacency(chunk.VoxelCoordinate{X: x, Y: y, Z: z}, voxInfo.adjMask)
				ch.SetBlockType(chunk.VoxelCoordinate{X: x, Y: y, Z: z}, voxInfo.bType)
			}
		}
	}
	return ch
}

func (gen *AlexWorldGenerator) generateAt(x, y, z int32) *genVoxel {
	vp := chunk.VoxelCoordinate{
		X: x,
		Y: y,
		Z: z,
	}
	var faceMods = [6]struct {
		off     chunk.VoxelCoordinate
		adjFace chunk.AdjacentMask
	}{
		{chunk.VoxelCoordinate{X: -1, Y: 0, Z: 0}, chunk.AdjacentLeft},
		{chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.AdjacentRight},
		{chunk.VoxelCoordinate{X: 0, Y: -1, Z: 0}, chunk.AdjacentBottom},
		{chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0}, chunk.AdjacentTop},
		{chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}, chunk.AdjacentFront},
		{chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1}, chunk.AdjacentBack},
	}
	var mask chunk.AdjacentMask
	for _, mod := range faceMods {
		offP := chunk.VoxelCoordinate{
			X: vp.X + mod.off.X,
			Y: vp.Y + mod.off.Y,
			Z: vp.Z + mod.off.Z,
		}
		offV := alexHelper(offP)
		if offV != chunk.BlockTypeAir {
			mask |= mod.adjFace
		}
	}
	return &genVoxel{
		adjMask: mask,
		bType:   alexHelper(vp),
	}
}

func alexHelper(pos chunk.VoxelCoordinate) chunk.BlockType {
	h := int(math.Round(noiseAt(int(pos.X), int(pos.Z))) + 10)
	if int(pos.Y) > h {
		return chunk.BlockTypeAir
	} else {
		return chunk.BlockTypeGrassSides
	}
}

const (
	rootTwo   = 1.4142135623730950488016887242096980785696718753769480731766797379
	rootThree = 1.7320508075688772935274463415058723669428052538103806280558069794
	rootSeven = 2.6457513110645905905016157536392604257102591830824501803683344592
)

func noiseAt(x, z int) float64 {
	res := 20.0
	return math.Round(10*(smoothNoise(float64(x)/res)+
		smoothNoise(float64(z)/res)+
		smoothNoise(float64(x+z)/res)+
		smoothNoise(float64(x-z)/res))) / 10
}

func roughNoise(x float64) float64 {
	return math.Cos(x)*(math.Sin(rootTwo*x)+math.Sin(rootThree*x)) + math.Sin(rootSeven*x)*(math.Cos(rootTwo*x)+math.Cos(rootThree*x))
}

func smoothNoise(x float64) float64 {
	r := 0.0
	samples := 10.0
	delta := 0.2
	for k := 0.0; k < samples; k++ {
		r += roughNoise(x + k*((2*delta)/samples))
	}

	return r / samples
}
