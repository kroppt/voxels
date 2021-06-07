package world_test

import (
	"math"
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

func TestStraightDownIntersect(t *testing.T) {
	t.Parallel()
	cam := world.NewCamera()
	cam.SetPosition(&glm.Vec3{10.5, 3.5, 10.5})
	cam.LookAt(&glm.Vec3{10.5, 0.5, 10.5})
	var root *world.Octree
	root = root.AddLeaf(&world.Voxel{
		Pos: world.VoxelPos{0, 0, 0},
	})
	root = root.AddLeaf(&world.Voxel{
		Pos: world.VoxelPos{10, 0, 10},
	})
	result, ok := root.Find(getCamIntersectionPredicate(cam))
	closest, _ := world.GetClosest(cam.GetPosition(), result)
	expectVoxel := world.VoxelPos{10, 0, 10}
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
}

func TestIntersectAlone(t *testing.T) {
	testCases := []struct {
		desc      string
		aabc      world.AABC
		expectHit bool
		pos       glm.Vec3
		target    glm.Vec3
	}{
		{
			desc: "simple straight down dead center",
			aabc: world.AABC{
				Origin: world.VoxelPos{0, 0, 0},
				Size:   1,
			},
			expectHit: true,
			pos:       glm.Vec3{0.5, 2, 0.5},
			target:    glm.Vec3{0.5, 1, 0.5},
		},
		{
			desc: "simple straight down close to edge",
			aabc: world.AABC{
				Origin: world.VoxelPos{0, 0, 0},
				Size:   1,
			},
			expectHit: true,
			pos:       glm.Vec3{0.99, 2, 0.5},
			target:    glm.Vec3{0.99, 1, 0.5},
		},
		{
			desc: "angle close to edge",
			aabc: world.AABC{
				Origin: world.VoxelPos{0, 0, 0},
				Size:   1,
			},
			expectHit: true,
			pos:       glm.Vec3{0.5, 2, 0.5},
			target:    glm.Vec3{0.99, 1, 0.5},
		},
		{
			desc: "other angle",
			aabc: world.AABC{
				Origin: world.VoxelPos{0, 0, 0},
				Size:   1,
			},
			expectHit: true,
			pos:       glm.Vec3{-4, 2, 0.77},
			target:    glm.Vec3{0.8, 1, 0.123123123},
		},
		{
			desc: "weird angle and far awar",
			aabc: world.AABC{
				Origin: world.VoxelPos{40, 0, 40},
				Size:   1,
			},
			expectHit: true,
			pos:       glm.Vec3{-44, 2, 40.77},
			target:    glm.Vec3{40.2, 1, 40.123123123},
		},
		{
			desc: "big aabc",
			aabc: world.AABC{
				Origin: world.VoxelPos{1, 1, 1},
				Size:   2,
			},
			expectHit: true,
			pos:       glm.Vec3{4, 6, 0},
			target:    glm.Vec3{2, 3, 2},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			cam := world.NewCamera()
			cam.SetPosition(&tC.pos)
			cam.LookAt(&tC.target)
			dist, hit := world.Intersect(tC.aabc, cam.GetPosition(), cam.GetLookForward())
			dx := float64(tC.pos.X() - tC.target.X())
			dy := float64(tC.pos.Y() - tC.target.Y())
			dz := float64(tC.pos.Z() - tC.target.Z())
			h1 := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
			expectDist := math.Sqrt(math.Pow(dz, 2) + math.Pow(h1, 2))
			allowableError := float32(0.0001)
			if hit != tC.expectHit {
				t.Fatalf("expected hit to be %v but got %v", tC.expectHit, hit)
			}
			if !withinError(float32(expectDist), dist, allowableError) {
				t.Fatalf("expected dist %v but got %v", expectDist, dist)
			}
		})
	}
}
