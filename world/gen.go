package world

type Generator interface {
	GenerateAt(int, int, int) *Voxel
}

type FlatWorldGenerator struct {
}

func (g1 FlatWorldGenerator) GenerateAt(x, y, z int) *Voxel {
	if y < 0 || y > 5 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll,
			Btype:   Air,
		}
	}
	if y == 0 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll & ^AdjacentBottom,
			Btype:   Labeled,
		}
	} else if y == 5 {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll & ^AdjacentTop,
			Btype:   Grass,
		}
	} else {
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll,
			Btype:   Dirt,
		}
	}
}
