package view_test

import (
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/view"
)

func TestOneVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   1,
	}
	resultAABC := root.GetAABC()
	if *resultAABC != *expectedAABC {
		t.Fatalf("expected AABC %v but got %v", *expectedAABC, *resultAABC)
	}
	expectedVoxel := &chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}
	resultVoxel := root.GetVoxel()
	if *resultVoxel != *expectedVoxel {
		t.Fatalf("expected Voxel %v but got %v", *expectedVoxel, *resultVoxel)
	}
	resultChildrenCount := root.CountChildren()
	expectedChildrenCount := 0
	if resultChildrenCount != expectedChildrenCount {
		t.Fatalf("expected %v children but counted %v", expectedChildrenCount, resultChildrenCount)
	}
}

func TestAdjacentTwoVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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
}

func TestCorneredTwoVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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
}

func TestTwoDistantVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: 0})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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

}

func TestTwoVeryDistantVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 4, Y: 0, Z: 0})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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
}

func TestTwoDistantVoxelOctreeWithAnother(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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

}

func TestOctreeReassignment(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	vOld := &chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}
	root = root.AddLeaf(vOld)
	vNew := &chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}
	root = root.AddLeaf(vNew)
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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

}

func TestThreeVoxelOctree(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
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

}

func TestThreeVoxelOctreeWithBackwards(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{-4, -4, -4},
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
}

func TestOctreeFarCornerDoesntChange(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -4, Y: 3, Z: 3})
	expectedAABC := &view.AABC{
		Origin: mgl.Vec3{-4, -4, -4},
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
}

func TestOctreeRemoveRoot(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})

	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result != nil {
		t.Fatalf("expected root to be removed, but wasn't")
	}
}

func TestOctreeDoNotRemoveRoot(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result == nil {
		t.Fatalf("root was removed when it shouldn't have been")
	}
}

func TestOctreeFastRootBreak(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 1, Z: 1})
	root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result != nil {
		t.Fatalf("expected root to be removed but wasn't")
	}
}

func TestOctreeFasterRootBreak(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result == nil {
		t.Fatalf("expected root to not be removed but was")
	}
}

func TestOctreeRootSingleShrink(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result.GetAABC().Size != 1 {
		t.Fatalf("expected root to be size 1 but was %v", result.GetAABC().Size)
	}
	if *(result.GetVoxel()) != (chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}) {
		t.Fatalf("expected pos to be {1,1,1} but got %v", *result.GetVoxel())
	}
}

func TestOctreeRootDoubleShrink(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result.GetAABC().Size != 2 {
		t.Fatalf("expected root to be size 2 but was %v", result.GetAABC().Size)
	}
}

func TestOctreeNoRootShrink(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	result, removed := root.Remove(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	if !removed {
		t.Fatalf("Expected a voxel to be removed, but root.Remove indicated that none were")
	}
	if result != root {
		t.Fatalf("expected the returned root from remove to be exactly the same, but wasn't")
	}
}

func TestFindClosestIntersectSingleVoxel(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3})
	eye := mgl.Vec3{3.5, 3.5, 4.5}
	dir := mgl.Vec3{0, 0, -1}
	expectDist := 0.5
	expectedVc := chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3}
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
}

func TestFindClosestIntersect(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: -2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: -3})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: -1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})

	eye := mgl.Vec3{0.5, 0.5, 0.5}
	dir := mgl.Vec3{0, 0, -1}
	expectDist := 0.5
	expectedVc := chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
}

func TestFindClosestIntersectFar(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: -3})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: -1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})

	eye := mgl.Vec3{0.5, 0.5, 0.5}
	dir := mgl.Vec3{0, 0, -1}
	expectDist := 2.5
	expectedVc := chunk.VoxelCoordinate{X: 0, Y: 0, Z: -3}
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
}

func TestFindClosestIntersectAngle(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: -2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: -1, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: -1, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})

	eye := mgl.Vec3{0.5, 0.5, 0.5}
	dir := mgl.Vec3{1, 0, -1}.Normalize()
	expectDist := 2.121
	expectedVc := chunk.VoxelCoordinate{X: 2, Y: 0, Z: -2}
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
}

func TestFindClosestIntersectAngle2(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 2, Y: 0, Z: -2})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: -2})

	eye := mgl.Vec3{0.5, 0.5, 0.5}
	dir := mgl.Vec3{1.5, 0, -2}.Normalize()
	expectedVc := chunk.VoxelCoordinate{X: 1, Y: 0, Z: -2}
	expectDist := 1.875
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
}

func TestFindClosestIntersectInsideVoxel(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 1})
	eye := mgl.Vec3{0.5, 0.5, 0.5}
	dir := mgl.Vec3{0, 0, -1}
	expectedVc := chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}
	expectDist := 0.0
	vc, dist, ok := root.FindClosestIntersect(eye, dir)
	if !ok {
		t.Fatal("expected to find intersect")
	}
	if vc != expectedVc {
		t.Fatalf("expected to intersect voxel %v but hit %v", expectedVc, vc)
	}
	if !withinError(dist, expectDist, errMargin) {
		t.Fatalf("expected dist %v but got %v (errMargin=%v)", expectDist, dist, errMargin)
	}
}

func TestTreePrint(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	actual := root.String()
	expected := "{(0,0,0)->{}}"
	if actual != expected {
		t.Fatalf("expected tree \"%v\" but got \"%v\"\n", expected, actual)
	}
}

func TestTreePrintTwoNodes(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})

	actual := root.String()
	expected := "{(2)->{(1,0,0)->{}}{(0,0,0)->{}}}"
	if actual != expected {
		t.Fatalf("expected tree \"%v\" but got \"%v\"\n", expected, actual)
	}
}

func TestTreePrintDoubleExpand(t *testing.T) {
	t.Parallel()
	var root *view.Octree
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	root = root.AddLeaf(&chunk.VoxelCoordinate{X: 3, Y: 0, Z: 0})

	actual := root.String()
	expected := "{(4)->{(2)->{(3,0,0)->{}}}{(2)->{(0,0,0)->{}}}}"
	if actual != expected {
		t.Fatalf("expected tree \"%v\" but got \"%v\"\n", expected, actual)
	}
}

func withinError(x, y float64, diff float64) bool {
	if x+diff > y && x-diff < y {
		return true
	}
	return false
}

const errMargin = 0.001
