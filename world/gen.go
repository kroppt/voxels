package world

import (
	"math"
)

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
	h := int(math.Round(noiseAt(pos.X, pos.Z, sectorAvgValue(pos.X, pos.Z))))

	if pos.Y == h {
		if h <= stoneLim {
			return Stone
		} else if h < snowLim {
			return Grass
		} else {
			return Light
		}
	} else if pos.Y < h {
		if pos.Y <= stoneLim {
			return Stone
		} else {
			return Dirt
		}
	} else {
		return Air
	}

	// sectors test
	// t := sectorValue(pos.X, pos.Z)
	// a := int(math.Round(sectorAvgValue(pos.X, pos.Z)))

	// if pos.Y > a {
	// 	return Air
	// } else if pos.Y == a && pos.X%sectorSize == 0 && pos.Z%sectorSize == 0 {
	// 	return Light
	// }

	// if t < 3 {
	// 	return Stone
	// } else if t < 6 {
	// 	return Dirt
	// } else {
	// 	return Grass
	// }
}

const (
	//DO NOT CHANGE THESE EVER!
	//changing these could result in terrain generation becoming periodic or uninteresting
	rootTwo   = 1.4142135623730950488016887242096980785696718753769480731766797379
	rootThree = 1.7320508075688772935274463415058723669428052538103806280558069794
	rootSeven = 2.6457513110645905905016157536392604257102591830824501803683344592

	//TERRAIN GENERATION CONSTANTS:
	//sectorSize being larger should result in larger flat/hilly areas, 20-30 looks good
	//I think sectorsize should be an even number, could have unintended results if not
	sectorSize = 30
	//max elevation a sector can have, 50 looks good even with a much lower snowlim
	//unlikely to be obtained since one random factor reduces this per sector, and another random factor achieves only part of a sector's max
	maxElevation = 50
	//res being larger will make the entire world smoother and more spread out, 10-15 looks good
	//low res values result in spiky behavior and the "trenches" are more evident
	res = 15
	//the seed will shift the noise function over, I think it will be a slightly different world and not just shifted
	//seed 0 should start camera by a mountain
	seed = 0
	//at and below stonelim, you get stone, 21 is good
	//at and above snowlim, you get "snow" (light blocks), 6 is good
	snowLim  = 21
	stoneLim = 6
)

func sectorAvgValue(x, z int) float64 {
	h := sectorSize / 2
	a := 0.0
	x1 := 0
	z1 := 0
	d := 0.0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			x1 = sectorSize*(div((i*sectorSize+x), sectorSize)) + h - 1
			z1 = sectorSize*(div((j*sectorSize+z), sectorSize)) + h - 1
			d = math.Min(distance(x, z, x1, z1), float64(sectorSize))
			a += (sectorSize - d) * sectorValue(x1, z1)
		}
	}
	return a / (sectorSize)
}

func distance(x, z, x1, z1 int) float64 {
	return math.Sqrt(float64((x-x1)*(x-x1) + (z-z1)*(z-z1)))
}

func sectorValue(x, z int) float64 {
	//uses the noise function for now, even though output is not uniform
	//the higher the number, the more elevation (makes it easy to blend by adding)
	return maxElevation * smoothNoise(float64(getSpiralIndex(div(x, sectorSize), div(z, sectorSize))))
}

func noiseAt(x, z int, max float64) float64 {
	return max * roundDigits(smoothNoise2D(float64(x)/res, float64(z)/res), 2)
}

func div(a, b int) int {
	//divides integers in the way they should be divided
	//not the bizzare way computer """scientists""" have decided they ought be rounded
	//assumes b is positive, if not I have no idea what happens
	if a >= 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

func roughNoise(x float64) float64 {
	x1 := x + seed
	r := math.Cos(x1)*(math.Sin(rootTwo*x1)+math.Sin(rootThree*x1)) + math.Sin(rootSeven*x1)*(math.Cos(rootTwo*x1)+math.Cos(rootThree*x1))
	// normalizes r to lie between 0 and 1
	r += 2.8258623603
	r /= 5.6524089098
	if r < 0 {
		r = 0
	} else if r > 1 {
		r = 1
	}
	return r
}

func smoothNoise(x float64) float64 {
	r := 0.0
	samples := 10.0
	delta := 0.1
	for k := 0.0; k < samples; k++ {
		r += roughNoise(x + k*((2*delta)/samples))
	}
	return r / samples
}

func smoothNoise2D(x, z float64) float64 {
	// a := 0.0
	// avgSize := 1
	// for i := -avgSize; i <= avgSize; i++ {
	// 	for j := -avgSize; j <= avgSize; j++ {
	// 		a += (smoothNoise((x + float64(i))) + smoothNoise(z+float64(j))) / 2
	// 	}
	// }
	// a /= float64((2*avgSize + 1) * (2*avgSize + 1))
	// return a

	// this makes trenches in x,z,x=z, and x=-z directions
	//tM := 3.0 //trench modifier
	// tO := 0.0 //trench offset
	//return (smoothNoise(x)*smoothNoise(tM*z) + smoothNoise(z)*smoothNoise(tM*x) + smoothNoise(x+z)*smoothNoise(tM*(x-z)) + smoothNoise(x-z)*smoothNoise(tM*(x+z))) / 4.0
	return (smoothNoise(x) + smoothNoise(z) + smoothNoise(x+(z/2)) + smoothNoise(x+(z)) + smoothNoise(x+(2*z)) + smoothNoise(x-(z/2)) + smoothNoise(x-(z)) + smoothNoise(x-(2*z))) / 8.0
}

func getSpiralIndex(x, z int) int {
	//returns index in spiral of x,y
	i := 0
	l := int(math.Max(math.Abs(float64(x)), math.Abs(float64(z))))
	if l == 0 {
		return 0
	}
	r := 2*l - 1 //r^2 already done
	//just trust me bro
	if x == l && z > -l {
		i = r*r + 0*l - 1 + l + z
	} else if z == l {
		i = r*r + 2*l - 1 + l - x
	} else if x == -l {
		i = r*r + 4*l - 1 + l - z
	} else if z == -l {
		i = r*r + 6*l - 1 + l + x
	}
	return i
}

func roundDigits(a float64, digits int) float64 {
	d := math.Pow(10, float64(digits))
	return math.Round(d*a) / d
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
