package world

import "math"

type Generator interface {
	GenerateAt(int, int, int) *Voxel
}

type AlexGenerator struct {
}

func (aa AlexGenerator) GenerateAt(x, y, z int) *Voxel {
	vp := VoxelPos{
		X: x,
		Y: y,
		Z: z,
	}
	var faceMods = [6]struct {
		off     VoxelPos
		adjFace AdjacentMask
	}{
		{VoxelPos{-1, 0, 0}, AdjacentLeft},
		{VoxelPos{1, 0, 0}, AdjacentRight},
		{VoxelPos{0, -1, 0}, AdjacentBottom},
		{VoxelPos{0, 1, 0}, AdjacentTop},
		{VoxelPos{0, 0, -1}, AdjacentFront},
		{VoxelPos{0, 0, 1}, AdjacentBack},
	}
	var mask AdjacentMask
	for _, mod := range faceMods {
		offP := vp.Add(mod.off)
		offV := alexHelper(offP)
		if offV != Air {
			mask |= mod.adjFace
		}
	}
	return &Voxel{
		Pos:     VoxelPos{x, y, z},
		AdjMask: mask,
		Btype:   alexHelper(vp),
	}
}

func alexHelper(pos VoxelPos) BlockType {
	h := int(math.Round(noiseAt(pos.X, pos.Z)) + 10)
	if pos.Y > h {
		return Air
	} else if pos.Y == h && pos.X%5 == 0 && pos.Z%5 == 0 {
		return Light
	} else {
		return Corrupted
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
