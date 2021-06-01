package world_test

import (
	"fmt"
	"testing"

	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
	"github.com/kroppt/voxels/world"
)

func TestTableAABBIntersect(t *testing.T) {

	testCases := []struct {
		campos     glm.Vec3
		camdir     glm.Vec3
		boxpos     glm.Vec3
		boxext     glm.Vec3
		intersects bool
	}{
		{
			campos:     glm.Vec3{-0, -0, 25},
			camdir:     glm.Vec3{0, 0, -1},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: true,
		},
		{
			campos:     glm.Vec3{-0, -0, 25},
			camdir:     glm.Vec3{-0.008726252, -0.0034906445, -0.99995583},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: true,
		},
		{
			campos:     glm.Vec3{1.5, -0, 25},
			camdir:     glm.Vec3{0, 0, -1},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 0, 25},
			camdir:     glm.Vec3{0, 0, 1},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 0, 25},
			camdir:     glm.Vec3{1, 0, 0},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 0, 25},
			camdir:     glm.Vec3{1, 0, 0},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 0, 25},
			camdir:     glm.Vec3{1, 0, 1},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 1, 0},
			camdir:     glm.Vec3{0, 0, 1},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: false,
		},
		{
			campos:     glm.Vec3{0, 1, 0},
			camdir:     glm.Vec3{-1, -1, 0},
			boxpos:     glm.Vec3{0, 0, 0},
			boxext:     glm.Vec3{0.5, 0.5, 0.5},
			intersects: true,
		},
	}

	for _, tC := range testCases {

		desc := fmt.Sprintf("campos %v camdir %v boxpos %v boxext %v", tC.campos, tC.camdir, tC.boxpos, tC.boxext)
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			cam := world.NewCamera()
			cam.SetPosition(&tC.campos)
			center := tC.campos.Add(&tC.camdir)
			cam.LookAt(&center)
			aabb := geo.AABB{
				Center:     tC.boxpos,
				HalfExtend: tC.boxext,
			}
			dist, hit := world.Intersect(aabb, cam.GetPosition(), cam.GetLookForward())
			if hit {
				t.Logf("dist: %v", dist)
			}
			if hit != tC.intersects {
				t.Fatalf("expected overlap=%v but got %v", tC.intersects, hit)
			}
		})

	}

}
