package world_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/world"
)

func getCamIntersectionPredicate(cam *world.Camera) func(*world.Octree) bool {
	return func(node *world.Octree) bool {
		aabc := *node.GetAABC()
		_, hit := world.Intersect(aabc, cam.GetPosition(), cam.GetLookForward())
		return hit
	}
}

func TestSimpleVoxelIntersect(t *testing.T) {
	t.Run("simple voxel-ray intersection", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		cam.SetPosition(&glm.Vec3{0.5, 0.5, -2})
		cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		result, ok := root.Find(getCamIntersectionPredicate(cam))
		closest, _ := world.GetClosest(cam.GetPosition(), result)
		expectVoxel := world.VoxelPos{0, 0, 0}
		expectedIntersections := 1
		if !ok {
			t.Fatal("view did not intersect voxel")
		}
		if len(result) != expectedIntersections {
			t.Fatal("only expected to intersect one voxel")
		}
		if closest.Pos != expectVoxel {
			t.Fatalf("expected to find voxel at %v but found %v", expectVoxel, closest.Pos)
		}

	})
}

func TestMultipleVoxelIntersect(t *testing.T) {
	t.Run("look through many and hit closer voxel", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		cam.SetPosition(&glm.Vec3{9, 0.5, 0.5})
		cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{4, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{3, 0, 7},
		})
		result, ok := root.Find(getCamIntersectionPredicate(cam))
		closest, _ := world.GetClosest(cam.GetPosition(), result)
		expectVoxel := world.VoxelPos{4, 0, 0}
		expectedIntersections := 3
		if !ok {
			t.Fatal("view did not intersect voxel")
		}
		if len(result) != expectedIntersections {
			t.Fatal("only expected to intersect one voxel")
		}
		if closest.Pos != expectVoxel {
			t.Fatalf("expected to find voxel at %v but found %v", expectVoxel, closest.Pos)
		}

	})
}

func TestMultipleVoxelIntersectLookBetween(t *testing.T) {
	t.Run("look between many voxels to see one in the back", func(t *testing.T) {
		t.Parallel()
		cam := world.NewCamera()
		cam.SetPosition(&glm.Vec3{3.5, 0.5, -10})
		cam.LookAt(&glm.Vec3{3.5, 0.5, 7.5})
		var root *world.Octree
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{0, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{2, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{4, 0, 0},
		})
		root = root.AddLeaf(&world.Voxel{
			Pos: world.VoxelPos{3, 0, 7},
		})
		result, ok := root.Find(getCamIntersectionPredicate(cam))
		closest, _ := world.GetClosest(cam.GetPosition(), result)
		expectVoxel := world.VoxelPos{3, 0, 7}
		expectedIntersections := 1
		if !ok {
			t.Fatal("view did not intersect voxel")
		}
		if len(result) != expectedIntersections {
			t.Fatal("only expected to intersect one voxel")
		}
		if closest.Pos != expectVoxel {
			t.Fatalf("expected to find voxel at %v but found %v", expectVoxel, closest.Pos)
		}

	})
}
