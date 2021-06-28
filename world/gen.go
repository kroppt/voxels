package world

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

type Generator interface {
	GenerateAt(int, int, int) *Voxel
}

type NoiseGenerator struct {
}

type FlatWorldGenerator struct {
}

func (g1 FlatWorldGenerator) GenerateAt(x, y, z int) *Voxel {
	if y < 0 || y > 6 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentNone,
			Btype:   Air,
		}
	}
	if y == 0 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll & ^AdjacentBottom,
			Btype:   Labeled,
		}
	} else if y == 6 {
		if x == 3 && z == 3 {
			return &Voxel{
				Pos:     VoxelPos{x, y, z},
				AdjMask: AdjacentAll & ^AdjacentTop,
				Btype:   Light,
			}
		} else {
			return &Voxel{
				Pos:     VoxelPos{x, y, z},
				AdjMask: AdjacentAll & ^AdjacentTop,
				Btype:   Grass,
			}
		}
	} else if y == 1 || y == 2 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll,
			Btype:   Corrupted,
		}
	} else if y == 3 || y == 4 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll,
			Btype:   Stone,
		}
	} else {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll,
			Btype:   Dirt,
		}
	}
}

const ( //paramters for noise generator
	seed       = 900000000
	lacunarity = 2
	gain       = 0.5
	octavec    = 5 //
	frequency  = 0.015
	amplitude  = 0.75
	arbitrary  = 50 // >=10 works best
	// TODO Figure out linear mapping so arbitrary is unnecessary
)

func Octaves(sd int64, frq, amp, oct, x, z, lac, gain, arb float32) float32 {
	seed := opensimplex.New32(sd)
	var sum float32
	for i := float32(0); i < oct; i++ {
		sum += amp * seed.Eval2(frq*x, frq*z)
		amp *= gain
		frq *= lac
	}
	return ((sum / oct) * arb)
}

func (g1 NoiseGenerator) GenerateAt(x, y, z int) *Voxel {
	// test := opensimplex.New32(seed)
	// value := test.Eval2(frequency*float32(x), frequency*float32(z))
	// value += amplitude * test.Eval2(frequency*2*float32(x), frequency*2*float32(z))
	// value += amplitude / 2 * test.Eval2(frequency*4*float32(x), frequency*4*float32(z))
	// value += amplitude / 4 * test.Eval2(frequency*8*float32(x), frequency*8*float32(z))
	// value += amplitude / 8 * test.Eval2(frequency*16*float32(x), frequency*16*float32(z))
	// avgVal := value / 5

	// // value += (0.5 * test.Eval2(.09*float32(x), .09*float32(z)))
	// avgVal *= arbitrary
	newY := int(math.Round(float64(Octaves(seed, frequency, amplitude, octavec, float32(x), float32(z), lacunarity, gain, arbitrary))))
	if y > newY {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentNone,
			Btype:   Air,
		}
	} else if y < 0 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentNone,
			Btype:   Corrupted,
		}
	} else {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentNone,
			Btype:   Grass,
		}
	}
}
