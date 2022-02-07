package world

import (
	"container/list"
	"math"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/repositories/settings"
)

type Generator interface {
	GenerateChunk(chunk.ChunkCoordinate) (chunk.Chunk, *list.List)
}

type FnGenerator struct {
	FnGenerateChunk func(chunk.ChunkCoordinate) (chunk.Chunk, *list.List)
}

func (fn *FnGenerator) GenerateChunk(pos chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
	if fn.FnGenerateChunk != nil {
		return fn.FnGenerateChunk(pos)
	}
	return chunk.NewChunkEmpty(pos, 1), list.New()
}

type TrentWorldGenerator struct {
	settingsRepo settings.Interface
}

func NewTrentWorldGenerator(settingsRepo settings.Interface) *TrentWorldGenerator {
	if settingsRepo == nil {
		panic("flat world generator missing settings repo")
	}
	return &TrentWorldGenerator{
		settingsRepo: settingsRepo,
	}
}

func (gen *TrentWorldGenerator) GenerateChunk(chPos chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
	size := int32(gen.settingsRepo.GetChunkSize())
	ch := chunk.NewChunkEmpty(chPos, uint32(size))
	pending := list.New()
	ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
		pending.PushBackList(ch.SetBlockType(vc, gen.generateAt(vc.X, vc.Y, vc.Z)))
	})
	return ch, pending
}

func (gen *TrentWorldGenerator) generateAt(x, y, z int32) chunk.BlockType {
	return chunk.BlockTypeAir
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

func (gen *FlatWorldGenerator) GenerateChunk(chPos chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
	size := int32(gen.settingsRepo.GetChunkSize())
	ch := chunk.NewChunkEmpty(chPos, uint32(size))
	pending := list.New()
	ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
		pending.PushBackList(ch.SetBlockType(vc, gen.generateAt(vc.X, vc.Y, vc.Z)))
	})
	return ch, pending
}

func (gen *FlatWorldGenerator) generateAt(x, y, z int32) chunk.BlockType {
	if y < 0 || y > 6 {
		return chunk.BlockTypeAir
	}
	if y == 0 {
		return chunk.BlockTypeLabeled
	} else if y == 6 {
		if x == 3 && z == 3 {
			return chunk.BlockTypeLight
		} else {
			return chunk.BlockTypeGrass
		}
	} else if y == 1 || y == 2 {
		return chunk.BlockTypeCorrupted

	} else if y == 3 || y == 4 {
		return chunk.BlockTypeStone
	} else {
		return chunk.BlockTypeDirt
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

func (gen *AlexWorldGenerator) GenerateChunk(chPos chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
	size := int32(gen.settingsRepo.GetChunkSize())
	ch := chunk.NewChunkEmpty(chPos, uint32(size))
	pending := list.New()
	ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
		pending.PushBackList(ch.SetBlockType(vc, alexHelper(vc)))
	})
	return ch, pending
}

func alexHelper(pos chunk.VoxelCoordinate) chunk.BlockType {
	h := int(math.Round(noiseAt(int(pos.X), int(pos.Z))) + 10)
	if int(pos.Y) > h {
		return chunk.BlockTypeAir
	} else if int(pos.Y) == h {
		return chunk.BlockTypeGrassSides
	} else if int(pos.Y) < h && int(pos.Y) > h-3 {
		return chunk.BlockTypeDirt
	} else {
		return chunk.BlockTypeStone
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
