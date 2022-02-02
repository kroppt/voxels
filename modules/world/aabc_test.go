package world_test

import (
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/modules/world"
)

func TestWithinAABCCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		aabc   world.AABC
		target mgl.Vec3
		expect bool
	}{
		{
			desc: "WithinAABC works for standard point",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC excludes point on maximum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 4, 4},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on X maximum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 1, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Y maximum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 4, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Z maximum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, 0},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on X minimum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{0, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Y minimum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 0, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Z minimum edge",
			aabc: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 0},
			expect: true,
		},
		{
			desc: "WithinAABC excludes far off point",
			aabc: world.AABC{
				Origin: mgl.Vec3{-1, -2, -3},
				Size:   2,
			},
			target: mgl.Vec3{-10, 10, 50},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point slightly outside",
			aabc: world.AABC{
				Origin: mgl.Vec3{-1, -2, -3},
				Size:   2,
			},
			target: mgl.Vec3{-1, -2, -4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge w/ negative center",
			aabc: world.AABC{
				Origin: mgl.Vec3{-5, -6, -7},
				Size:   4,
			},
			target: mgl.Vec3{-5, -6, -7},
			expect: true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.WithinAABC(tC.aabc, tC.target)
			if result != tC.expect {
				t.Fatalf("got %v but expected %v", result, tC.expect)
			}
		})
	}
}

func TestExpandAABCInsidePanic(t *testing.T) {
	t.Parallel()
	defer func() {
		recover()
	}()
	curr := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{1, 1, 1}
	world.ExpandAABC(curr, target)
	t.Fatal("expected panic but did not")
}

func TestExpandAABCOutsideDoesntPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("expected no panic, but did anyway")
		}
	}()
	curr := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{5, 6, 7}
	world.ExpandAABC(curr, target)
}

func TestExpandAABCCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		curr   world.AABC
		target mgl.Vec3
		expect world.AABC
	}{
		{
			desc: "ExpandAABC is bigger than current",
			curr: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{5, 5, 5},
			expect: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC maximum is exclusive",
			curr: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 4, 4},
			expect: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC minimum is inclusive",
			curr: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
			target: mgl.Vec3{-8, -8, -8},
			expect: world.AABC{
				Origin: mgl.Vec3{-8, -8, -8},
				Size:   16,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.ExpandAABC(tC.curr, tC.target)
			if result != tC.expect {
				t.Fatalf("got AABC %v but expected %v", result, tC.expect)
			}
		})
	}
}

func TestExpandAABCDoublesSize(t *testing.T) {
	t.Parallel()
	curr := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{9, 9, 9}
	expect := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   16,
	}
	step1 := world.ExpandAABC(curr, target)
	result := world.ExpandAABC(step1, target)
	if result != expect {
		t.Fatalf("got AABC %v but expected %v", result, expect)
	}
}

func TestGetChildAABCOutsidePanic(t *testing.T) {
	t.Parallel()
	defer func() {
		recover()
	}()
	aabc := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{-1, -1, -1}
	world.GetChildAABC(aabc, target)
	t.Fatal("expected panic but did not")
}

func TestGetChildAABCReturnsOctant(t *testing.T) {
	t.Parallel()
	defer func() {
		recover()
	}()
	aabc := world.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{-1, -1, -1}
	result := world.GetChildAABC(aabc, target)
	volume := result.Size * result.Size * result.Size
	expect := 8
	if volume != expect {
		t.Fatalf("expected octant volume to be %v but got %v", expect, volume)
	}
}

func TestGetChildAABCCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		curr   world.AABC
		target mgl.Vec3
		expect world.AABC
	}{
		{
			desc: "GetChildAABC +x +y +z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, 0},
			expect: world.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x +y -z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, -1},
			expect: world.AABC{
				Origin: mgl.Vec3{0, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y +z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, -1, 0},
			expect: world.AABC{
				Origin: mgl.Vec3{0, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y -z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, -1, -1},
			expect: world.AABC{
				Origin: mgl.Vec3{0, -2, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y +z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, 0, 0},
			expect: world.AABC{
				Origin: mgl.Vec3{-2, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y -z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, 0, -1},
			expect: world.AABC{
				Origin: mgl.Vec3{-2, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y +z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, -1, 0},
			expect: world.AABC{
				Origin: mgl.Vec3{-2, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y -z",
			curr: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, -1, -1},
			expect: world.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   2,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := world.GetChildAABC(tC.curr, tC.target)
			if result != tC.expect {
				t.Fatalf("got AABC %v but expected %v", result, tC.expect)
			}
		})
	}
}
