package oldworld_test

import (
	"testing"

	oldworld "github.com/kroppt/voxels/oldworld"
)

func TestWithinAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		aabc   *oldworld.AABC
		target oldworld.VoxelPos
		expect bool
	}{
		{
			desc: "WithinAABC works for standard point",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{1, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC excludes point on maximum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{4, 4, 4},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on X maximum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{4, 1, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Y maximum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{1, 4, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Z maximum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{1, 1, 4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, 0, 0},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on X minimum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Y minimum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{1, 0, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Z minimum edge",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{1, 1, 0},
			expect: true,
		},
		{
			desc: "WithinAABC excludes far off point",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-1, -2, -3},
				Size:   2,
			},
			target: oldworld.VoxelPos{-10, 10, 50},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point slightly outside",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-1, -2, -3},
				Size:   2,
			},
			target: oldworld.VoxelPos{-1, -2, -4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge w/ negative center",
			aabc: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-5, -6, -7},
				Size:   4,
			},
			target: oldworld.VoxelPos{-5, -6, -7},
			expect: true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := oldworld.WithinAABC(tC.aabc, tC.target)
			if result != tC.expect {
				t.Fatalf("got %v but expected %v", result, tC.expect)
			}
		})
	}
}

func TestExpandAABCInsidePanic(t *testing.T) {
	t.Run("ExpandAABC panics if target is within", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		curr := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   4,
		}
		target := oldworld.VoxelPos{1, 1, 1}
		oldworld.ExpandAABC(curr, target)
		t.Fatal("expected panic but did not")
	})
}

func TestExpandAABCOutsideDoesntPanic(t *testing.T) {
	t.Run("ExpandAABC does not panic normally", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err != nil {
				t.Fatal("expected no panic, but did anyway")
			}
		}()
		curr := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   4,
		}
		target := oldworld.VoxelPos{5, 6, 7}
		oldworld.ExpandAABC(curr, target)
	})
}

func TestExpandAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *oldworld.AABC
		target oldworld.VoxelPos
		expect *oldworld.AABC
	}{
		{
			desc: "ExpandAABC is bigger than current",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{5, 5, 5},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC maximum is exclusive",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   4,
			},
			target: oldworld.VoxelPos{4, 4, 4},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   8,
			},
		},
		{
			desc: "ExpandAABC minimum is inclusive",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   8,
			},
			target: oldworld.VoxelPos{-8, -8, -8},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-8, -8, -8},
				Size:   16,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := oldworld.ExpandAABC(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABC %v but expected %v", *result, *tC.expect)
			}
		})
	}
}

func TestExpandAABCDoublesSize(t *testing.T) {
	t.Run("ExpandAABC increases size twice", func(t *testing.T) {
		t.Parallel()
		curr := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   4,
		}
		target := oldworld.VoxelPos{9, 9, 9}
		expect := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   16,
		}
		step1 := oldworld.ExpandAABC(curr, target)
		result := oldworld.ExpandAABC(step1, target)
		if *result != *expect {
			t.Fatalf("got AABC %v but expected %v", *result, *expect)
		}
	})
}

func TestGetChildAABCOutsidePanic(t *testing.T) {
	t.Run("GetChildAABC panics if target is not within", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		aabc := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   4,
		}
		target := oldworld.VoxelPos{-1, -1, -1}
		oldworld.GetChildAABC(aabc, target)
		t.Fatal("expected panic but did not")
	})
}

func TestGetChildAABCReturnsOctant(t *testing.T) {
	t.Run("GetChildAABC returns 1/8th of the volume", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		aabc := &oldworld.AABC{
			Origin: oldworld.VoxelPos{0, 0, 0},
			Size:   4,
		}
		target := oldworld.VoxelPos{-1, -1, -1}
		result := oldworld.GetChildAABC(aabc, target)
		volume := result.Size * result.Size * result.Size
		expect := 8
		if volume != expect {
			t.Fatalf("expected octant volume to be %v but got %v", expect, volume)
		}
	})
}

func TestGetChildAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *oldworld.AABC
		target oldworld.VoxelPos
		expect *oldworld.AABC
	}{
		{
			desc: "GetChildAABC +x +y +z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, 0, 0},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x +y -z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, 0, -1},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y +z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, -1, 0},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC +x -y -z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{0, -1, -1},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{0, -2, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y +z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{-1, 0, 0},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, 0, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x +y -z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{-1, 0, -1},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, 0, -2},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y +z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{-1, -1, 0},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, 0},
				Size:   2,
			},
		},
		{
			desc: "GetChildAABC -x -y -z",
			curr: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   4,
			},
			target: oldworld.VoxelPos{-1, -1, -1},
			expect: &oldworld.AABC{
				Origin: oldworld.VoxelPos{-2, -2, -2},
				Size:   2,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := oldworld.GetChildAABC(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABC %v but expected %v", *result, *tC.expect)
			}
		})
	}
}
