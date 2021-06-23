package world

type Generator interface {
	GenerateAt(int, int, int) *Voxel
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
		return &Voxel{
			Pos:     VoxelPos{x, y, z},
			AdjMask: AdjacentAll & ^AdjacentTop,
			Btype:   Grass,
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
