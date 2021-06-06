package world_test

import (
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestOneVoxelOctree(t *testing.T) {
	t.Run("build a tree with only 1 voxel", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   1,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}
		expectedVoxel := &world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		}
		resultVoxel := root.GetVoxel()
		if *resultVoxel != *expectedVoxel {
			t.Fatalf("expected Voxel %v but got %v", *expectedVoxel, *resultVoxel)
		}
		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 0
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})
}

func TestAdjacentTwoVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   2,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})
}

func TestCorneredTwoVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 kitty-corner voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   2,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})
}

func TestTwoDistantVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 distant voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   4,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})

}

func TestTwoVeryDistantVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 distant voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{4, 0, 0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   8,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})

}
func TestTwoDistantVoxelOctreeWithAnother(t *testing.T) {
	t.Run("make a tree with 2 distance voxels plus a close one", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   4,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})

}

func TestOctreeReassignment(t *testing.T) {
	t.Run("adding duplicate voxel in same position should overwrite", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{0, 0, 0},
			Color: world.Color{0.5, 0.5, 0.5, 0.5},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{0, 0, 0},
			Color: world.Color{1.0, 1.0, 1.0, 1.0},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   1,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 0
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}

		resultCol := root.GetVoxel().Color
		expectCol := world.Color{1.0, 1.0, 1.0, 1.0}
		if expectCol != resultCol {
			t.Fatalf("expected reassigned voxel color to be %v but got %v", expectCol, resultCol)
		}
	})

}

func TestOctreeRecursionReassignment(t *testing.T) {
	t.Run("adding duplicate voxel in same position should overwrite", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{0, 0, 0},
			Color: world.Color{0.5, 0.5, 0.5, 0.5},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{4, 0, 0},
			Color: world.Color{0.2, 0.5, 0.5, 0.5},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{2, 0, 0},
			Color: world.Color{0.5, 0.5, 0.5, 0.4},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos:   world.VoxelPos{2, 0, 0},
			Color: world.Color{1.0, 1.0, 1.0, 1.0},
		})
		list, _ := root.Find(func(o *world.Octree) bool {
			return true
		})

		var answer *world.Voxel
		for _, v := range list {
			if v.Pos.X == 2 {
				answer = v
			}
		}

		resultCol := answer.Color
		expectCol := world.Color{1.0, 1.0, 1.0, 1.0}
		if expectCol != resultCol {
			t.Fatalf("expected reassigned voxel color to be %v but got %v", expectCol, resultCol)
		}
	})

}

func TestThreeVoxelOctree(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 2, 2},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{0, 0, 0},
			Size:   4,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 3
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})

}
func TestThreeVoxelOctreeWithBackwards(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 2, 2},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{-1, -1, -1},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{-4, -4, -4},
			Size:   8,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})
}
func TestOctreeFarCornerDoesntChange(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 2, 2},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{-1, -1, -1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{-4, 3, 3},
		})
		expectedAABC := &world.AABC{
			Origin: world.VoxelPos{-4, -4, -4},
			Size:   8,
		}
		resultAABC := root.GetAABC()
		if *resultAABC != *expectedAABC {
			t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 3
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
		}
	})
}
