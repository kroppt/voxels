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
		vOld := &world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		}
		root = root.AddLeaf(vOld)
		vNew := &world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		}
		root = root.AddLeaf(vNew)
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

		resultVox := root.GetVoxel()
		if vOld == resultVox {
			t.Fatal("expected old voxel to be reassigned but it wasn't")
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

func TestOctreeRemoveRoot(t *testing.T) {
	t.Run("fill 2x2 tree and then remove all", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		root.Remove(world.VoxelPos{0, 0, 0})
		root.Remove(world.VoxelPos{1, 0, 0})
		root.Remove(world.VoxelPos{0, 1, 0})
		root.Remove(world.VoxelPos{1, 1, 0})
		root.Remove(world.VoxelPos{0, 0, 1})
		root.Remove(world.VoxelPos{1, 0, 1})
		root.Remove(world.VoxelPos{0, 1, 1})
		result, removed := root.Remove(world.VoxelPos{1, 1, 1})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result != nil {
			t.Fatalf("expected root to be removed, but wasn't")
		}
	})
}

func TestOctreeDoNotRemoveRoot(t *testing.T) {
	t.Run("fill 2x2 tree and then remove some, root preserved", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		root.Remove(world.VoxelPos{0, 0, 0})
		root.Remove(world.VoxelPos{0, 1, 0})
		root.Remove(world.VoxelPos{0, 1, 1})
		root.Remove(world.VoxelPos{0, 0, 1})
		result, removed := root.Remove(world.VoxelPos{1, 1, 1})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result == nil {
			t.Fatalf("root was removed when it shouldn't have been")
		}
	})
}

func TestOctreeFastRootBreak(t *testing.T) {
	t.Run("fast way to break remove logic", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		root.Remove(world.VoxelPos{0, 0, 0})
		root.Remove(world.VoxelPos{0, 1, 0})
		root.Remove(world.VoxelPos{0, 1, 1})
		root.Remove(world.VoxelPos{0, 0, 1})
		result, removed := root.Remove(world.VoxelPos{1, 1, 1})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result != nil {
			t.Fatalf("expected root to be removed but wasn't")
		}
	})
}

func TestOctreeFasterRootBreak(t *testing.T) {
	t.Run("faster way to break remove logic", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 1, 0},
		})
		result, removed := root.Remove(world.VoxelPos{0, 0, 0})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result == nil {
			t.Fatalf("expected root to not be removed but was")
		}
	})
}

func TestOctreeRootSingleShrink(t *testing.T) {
	t.Run("single shrink", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		result, removed := root.Remove(world.VoxelPos{0, 0, 0})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result.GetAABC().Size != 1 {
			t.Fatalf("expected root to be size 1 but was %v", result.GetAABC().Size)
		}
		if result.GetVoxel().Pos != (world.VoxelPos{1, 1, 1}) {
			t.Fatalf("expected pos to be {1,1,1} but got %v", result.GetVoxel().Pos)
		}
	})
}

func TestOctreeRootDoubleShrink(t *testing.T) {
	t.Run("double shrink", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 2, 2},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{3, 3, 3},
		})
		result, removed := root.Remove(world.VoxelPos{0, 0, 0})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result.GetAABC().Size != 2 {
			t.Fatalf("expected root to be size 2 but was %v", result.GetAABC().Size)
		}
	})
}

func TestOctreeNoRootShrink(t *testing.T) {
	t.Run("no shrink", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 1, 1},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{1, 0, 0},
		})
		result, removed := root.Remove(world.VoxelPos{0, 0, 0})
		if !removed {
			t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
		}
		if result != root {
			t.Fatalf("expected the returned root from remove to be exactly the same, but wasn't")
		}
	})
}
