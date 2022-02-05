package view_test

import (
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/modules/view"
)

func TestWithinAABCCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		aabc   view.AABC
		target mgl.Vec3
		expect bool
	}{
		{
			desc: "WithinAABC works for standard point",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC excludes point on maximum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 4, 4},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on X maximum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 1, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Y maximum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 4, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Z maximum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, 0},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on X minimum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{0, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Y minimum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 0, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Z minimum edge",
			aabc: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{1, 1, 0},
			expect: true,
		},
		{
			desc: "WithinAABC excludes far off point",
			aabc: view.AABC{
				Origin: mgl.Vec3{-1, -2, -3},
				Size:   2,
			},
			target: mgl.Vec3{-10, 10, 50},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point slightly outside",
			aabc: view.AABC{
				Origin: mgl.Vec3{-1, -2, -3},
				Size:   2,
			},
			target: mgl.Vec3{-1, -2, -4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge w/ negative center",
			aabc: view.AABC{
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
			result := view.WithinAABC(tC.aabc, tC.target)
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
	curr := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{1, 1, 1}
	view.ExpandAABC(curr, target)
	t.Fatal("expected panic but did not")
}

func TestExpandAABCOutsideDoesntPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("expected no panic, but did anyway")
		}
	}()
	curr := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{5, 6, 7}
	view.ExpandAABC(curr, target)
}

func TestExpandAABCCases(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		curr   view.AABC
		target mgl.Vec3
		expect view.AABC
	}{
		{
			desc: "ExpandAABC is bigger than current",
			curr: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{5, 5, 5},
			expect: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC maximum is exclusive",
			curr: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   4,
			},
			target: mgl.Vec3{4, 4, 4},
			expect: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC minimum is inclusive",
			curr: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   8,
			},
			target: mgl.Vec3{-8, -8, -8},
			expect: view.AABC{
				Origin: mgl.Vec3{-8, -8, -8},
				Size:   16,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := view.ExpandAABC(tC.curr, tC.target)
			if result != tC.expect {
				t.Fatalf("got AABC %v but expected %v", result, tC.expect)
			}
		})
	}
}

func TestExpandAABCDoublesSize(t *testing.T) {
	t.Parallel()
	curr := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{9, 9, 9}
	expect := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   16,
	}
	step1 := view.ExpandAABC(curr, target)
	result := view.ExpandAABC(step1, target)
	if result != expect {
		t.Fatalf("got AABC %v but expected %v", result, expect)
	}
}

func TestGetChildAABCOutsidePanic(t *testing.T) {
	t.Parallel()
	defer func() {
		recover()
	}()
	aabc := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{-1, -1, -1}
	view.GetChildAABC(aabc, target)
	t.Fatal("expected panic but did not")
}

func TestGetChildAABCReturnsOctant(t *testing.T) {
	t.Parallel()
	defer func() {
		recover()
	}()
	aabc := view.AABC{
		Origin: mgl.Vec3{0, 0, 0},
		Size:   4,
	}
	target := mgl.Vec3{-1, -1, -1}
	result := view.GetChildAABC(aabc, target)
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
		curr   view.AABC
		target mgl.Vec3
		expect view.AABC
	}{
		{
			desc: "GetChildAABC +x +y +z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, 0},
			expect: view.AABC{
				Origin: mgl.Vec3{0, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x +y -z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, 0, -1},
			expect: view.AABC{
				Origin: mgl.Vec3{0, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y +z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, -1, 0},
			expect: view.AABC{
				Origin: mgl.Vec3{0, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y -z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{0, -1, -1},
			expect: view.AABC{
				Origin: mgl.Vec3{0, -2, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y +z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, 0, 0},
			expect: view.AABC{
				Origin: mgl.Vec3{-2, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y -z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, 0, -1},
			expect: view.AABC{
				Origin: mgl.Vec3{-2, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y +z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, -1, 0},
			expect: view.AABC{
				Origin: mgl.Vec3{-2, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y -z",
			curr: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   4,
			},
			target: mgl.Vec3{-1, -1, -1},
			expect: view.AABC{
				Origin: mgl.Vec3{-2, -2, -2},
				Size:   2,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := view.GetChildAABC(tC.curr, tC.target)
			if result != tC.expect {
				t.Fatalf("got AABC %v but expected %v", result, tC.expect)
			}
		})
	}
}

var Dist float64
var Hit bool

func BenchmarkIntersect(b *testing.B) {
	aabc := view.AABC{
		Origin: [3]float64{0, 0, 0},
		Size:   1,
	}
	pos := mgl.Vec3{-0.5, 0.5, 1.5}
	dir := mgl.Vec3{1, 0, -1}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Dist, Hit = view.Intersect(aabc, pos, dir)
	}
}
