package world_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
	"github.com/kroppt/voxels/world"
)

func TestWithinAABBCases(t *testing.T) {
	testCases := []struct {
		desc   string
		aabb   *geo.AABB
		target glm.Vec3
		expect bool
	}{
		{
			desc: "WithinAABB works for standard point",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{1, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABB excludes point on maximum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{4, 4, 4},
			expect: false,
		},
		{
			desc: "WithinAABB excludes point on X maximum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{4, 1, 1},
			expect: false,
		},
		{
			desc: "WithinAABB excludes point on Y maximum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{1, 4, 1},
			expect: false,
		},
		{
			desc: "WithinAABB excludes point on Z maximum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{1, 1, 4},
			expect: false,
		},
		{
			desc: "WithinAABB includes point on minimum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, 0, 0},
			expect: true,
		},
		{
			desc: "WithinAABB includes point on X minimum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABB includes point on Y minimum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{1, 0, 1},
			expect: true,
		},
		{
			desc: "WithinAABB includes point on Z minimum edge",
			aabb: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{1, 1, 0},
			expect: true,
		},
		{
			desc: "WithinAABB excludes far off point",
			aabb: &geo.AABB{
				Center:     glm.Vec3{-1, -2, -3},
				HalfExtend: glm.Vec3{2, 5, 1},
			},
			target: glm.Vec3{-10, 10, 50},
			expect: false,
		},
		{
			desc: "WithinAABB excludes point slightly outside",
			aabb: &geo.AABB{
				Center:     glm.Vec3{-1, -2, -3},
				HalfExtend: glm.Vec3{2, 5, 1},
			},
			target: glm.Vec3{-1, -2, -4},
			expect: false,
		},
		{
			desc: "WithinAABB includes point on minimum edge w/ negative center",
			aabb: &geo.AABB{
				Center:     glm.Vec3{-5, -6, -7},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{-5, -6, -7},
			expect: true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.WithinAABB(tC.aabb, tC.target)
			if result != tC.expect {
				t.Fatalf("got %v but expected %v", result, tC.expect)
			}
		})
	}
}

func TestExpandAABBInsidePanic(t *testing.T) {
	t.Run("ExpandAABB panics if target is within", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		curr := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{2, 2, 2},
		}
		target := glm.Vec3{1, 1, 1}
		world.ExpandAABB(curr, target)
		t.Fatal("expected panic but did not")
	})
}

func TestExpandAABBOutsideDoesntPanic(t *testing.T) {
	t.Run("ExpandAABB does not panic normally", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err != nil {
				t.Fatal("expected no panic, but did anyway")
			}
		}()
		curr := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{2, 2, 2},
		}
		target := glm.Vec3{5, 6, 7}
		world.ExpandAABB(curr, target)
	})
}

func TestExpandAABBCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *geo.AABB
		target glm.Vec3
		expect *geo.AABB
	}{
		{
			desc: "ExpandAABB is bigger than current",
			curr: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{5, 5, 5},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{4, 4, 4},
			},
		},
		{
			desc: "ExpandAABB maximum is exclusive",
			curr: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{4, 4, 4},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{4, 4, 4},
			},
		},
		{
			desc: "ExpandAABB minimum is inclusive",
			curr: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{4, 4, 4},
			},
			target: glm.Vec3{-8, -8, -8},
			expect: &geo.AABB{
				Center:     glm.Vec3{-8, -8, -8},
				HalfExtend: glm.Vec3{8, 8, 8},
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.ExpandAABB(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABB %v but expected %v", *result, *tC.expect)
			}
		})
	}
}

func TestExpandAABBDoublesSize(t *testing.T) {
	t.Run("ExpandAABB increases size twice", func(t *testing.T) {
		t.Parallel()
		curr := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{2, 2, 2},
		}
		target := glm.Vec3{9, 9, 9}
		expect := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{8, 8, 8},
		}
		step1 := world.ExpandAABB(curr, target)
		result := world.ExpandAABB(step1, target)
		if *result != *expect {
			t.Fatalf("got AABB %v but expected %v", *result, *expect)
		}
	})
}

func TestGetChildAABBOutsidePanic(t *testing.T) {
	t.Run("GetChildAABB panics if target is not within", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		aabb := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{2, 2, 2},
		}
		target := glm.Vec3{-1, -1, -1}
		world.GetChildAABB(aabb, target)
		t.Fatal("expected panic but did not")
	})
}

func TestGetChildAABBReturnsOctant(t *testing.T) {
	t.Run("GetChildAABB returns 1/8th of the volume", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		aabb := &geo.AABB{
			Center:     glm.Vec3{0, 0, 0},
			HalfExtend: glm.Vec3{2, 2, 2},
		}
		target := glm.Vec3{-1, -1, -1}
		result := world.GetChildAABB(aabb, target)
		volume := result.HalfExtend.X() * result.HalfExtend.X() * result.HalfExtend.X()
		expect := float32(8.0)
		if volume != expect {
			t.Fatalf("expected octant volume to be %v but got %v", expect, volume)
		}
	})
}

func TestGetChildAABBCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *geo.AABB
		target glm.Vec3
		expect *geo.AABB
	}{
		{
			desc: "GetChildAABB +x +y +z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, 0, 0},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, 0, 0},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB +x +y -z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, 0, -1},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, 0, -2},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB +x -y +z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, -1, 0},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, -2, 0},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB +x -y -z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{0, -1, -1},
			expect: &geo.AABB{
				Center:     glm.Vec3{0, -2, -2},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB -x +y +z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{-1, 0, 0},
			expect: &geo.AABB{
				Center:     glm.Vec3{-2, 0, 0},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB -x +y -z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{-1, 0, -1},
			expect: &geo.AABB{
				Center:     glm.Vec3{-2, 0, -2},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB -x -y +z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{-1, -1, 0},
			expect: &geo.AABB{
				Center:     glm.Vec3{-2, -2, 0},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
		{
			desc: "GetChildAABB -x -y -z",
			curr: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{2, 2, 2},
			},
			target: glm.Vec3{-1, -1, -1},
			expect: &geo.AABB{
				Center:     glm.Vec3{-2, -2, -2},
				HalfExtend: glm.Vec3{1, 1, 1},
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.GetChildAABB(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABB %v but expected %v", *result, *tC.expect)
			}
		})
	}
}
