package world_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
	"github.com/kroppt/voxels/world"
)

func TestOneVoxelOctree(t *testing.T) {
	t.Run("build a tree with only 1 voxel", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{0.5, 0.5, 0.5},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}
		expectedVoxel := &world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		}
		resultVoxel := root.GetVoxel()
		if *resultVoxel != *expectedVoxel {
			t.Fatalf("expected Voxel %v but got %v", *expectedVoxel, *resultVoxel)
		}
		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 0
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})
}

func TestAdjacentTwoVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{1, 0, 0},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{1, 1, 1},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})
}

func TestCorneredTwoVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 kitty-corner voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{1, 1, 1},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{1, 1, 1},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})
}

func TestTwoDistantVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 distant voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 0, 0},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{2, 2, 2},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})

}

func TestTwoVeryDistantVoxelOctree(t *testing.T) {
	t.Run("make a tree with 2 distant voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{4, 0, 0},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{4, 4, 4},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})

}
func TestTwoDistantVoxelOctreeWithAnother(t *testing.T) {
	t.Run("make a tree with 2 distance voxels plus a close one", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{1, 0, 0},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{2, 2, 2},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})

}

func TestThreeVoxelOctree(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 2, 2},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{0, 0, 0},
			HalfExtend: [3]float32{2, 2, 2},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 3
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})

}
func TestThreeVoxelOctreeWithBackwards(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 2, 2},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{-1, -1, -1},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{-4, -4, -4},
			HalfExtend: [3]float32{4, 4, 4},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 2
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})
}
func TestOctreeFarCornerDoesntChange(t *testing.T) {
	t.Run("make a tree with 3 adjacent voxels", func(t *testing.T) {
		t.Parallel()
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{2, 2, 2},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{-1, -1, -1},
		})
		root = root.AddLeaf(&world.Voxel{
			Object: nil,
			Pos:    glm.Vec3{-4, 3, 3},
		})
		expectedAABB := &geo.AABB{
			Center:     [3]float32{-4, -4, -4},
			HalfExtend: [3]float32{4, 4, 4},
		}
		resultAABB := root.GetAABB()
		if *resultAABB != *expectedAABB {
			t.Fatalf("expected AABB %v but got %v", *expectedAABB, *resultAABB)
		}

		resultChildrenCount := root.CountChildren()
		expectedChildrenCount := 3
		if resultChildrenCount != expectedChildrenCount {
			t.Fatalf("expected %v children but counted %v", resultChildrenCount, expectedChildrenCount)
		}
	})
}
