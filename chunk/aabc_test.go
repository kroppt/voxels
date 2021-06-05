package chunk_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/chunk"
)

func TestWithinAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		aabc   *chunk.AABC
		target glm.Vec3
		expect bool
	}{
		{
			desc: "WithinAABC works for standard point",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{1, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC excludes point on maximum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{4, 4, 4},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on X maximum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{4, 1, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Y maximum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{1, 4, 1},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point on Z maximum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{1, 1, 4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{0, 0, 0},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on X minimum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{0, 1, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Y minimum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{1, 0, 1},
			expect: true,
		},
		{
			desc: "WithinAABC includes point on Z minimum edge",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{1, 1, 0},
			expect: true,
		},
		{
			desc: "WithinAABC excludes far off point",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{-1, -2, -3},
				Size: 2,
			},
			target: glm.Vec3{-10, 10, 50},
			expect: false,
		},
		{
			desc: "WithinAABC excludes point slightly outside",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{-1, -2, -3},
				Size: 2,
			},
			target: glm.Vec3{-1, -2, -4},
			expect: false,
		},
		{
			desc: "WithinAABC includes point on minimum edge w/ negative center",
			aabc: &chunk.AABC{
				Pos:  glm.Vec3{-5, -6, -7},
				Size: 4,
			},
			target: glm.Vec3{-5, -6, -7},
			expect: true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := chunk.WithinAABC(tC.aabc, tC.target)
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
		curr := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 4,
		}
		target := glm.Vec3{1, 1, 1}
		chunk.ExpandAABC(curr, target)
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
		curr := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 4,
		}
		target := glm.Vec3{5, 6, 7}
		chunk.ExpandAABC(curr, target)
	})
}

func TestExpandAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *chunk.AABC
		target glm.Vec3
		expect *chunk.AABC
	}{
		{
			desc: "ExpandAABC is bigger than current",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{5, 5, 5},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 8,
			},
		},
		{
			desc: "ExpandAABC maximum is exclusive",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 4,
			},
			target: glm.Vec3{4, 4, 4},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 8,
			},
		},
		{
			desc: "ExpandAABC minimum is inclusive",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 8,
			},
			target: glm.Vec3{-8, -8, -8},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{-8, -8, -8},
				Size: 16,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := chunk.ExpandAABC(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABC %v but expected %v", *result, *tC.expect)
			}
		})
	}
}

func TestExpandAABCDoublesSize(t *testing.T) {
	t.Run("ExpandAABC increases size twice", func(t *testing.T) {
		t.Parallel()
		curr := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 4,
		}
		target := glm.Vec3{9, 9, 9}
		expect := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 16,
		}
		step1 := chunk.ExpandAABC(curr, target)
		result := chunk.ExpandAABC(step1, target)
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
		aabc := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 4,
		}
		target := glm.Vec3{-1, -1, -1}
		chunk.GetChildAABC(aabc, target)
		t.Fatal("expected panic but did not")
	})
}

func TestGetChildAABCReturnsOctant(t *testing.T) {
	t.Run("GetChildAABC returns 1/8th of the volume", func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		aabc := &chunk.AABC{
			Pos:  glm.Vec3{0, 0, 0},
			Size: 4,
		}
		target := glm.Vec3{-1, -1, -1}
		result := chunk.GetChildAABC(aabc, target)
		volume := result.Size * result.Size * result.Size
		expect := float32(8.0)
		if volume != expect {
			t.Fatalf("expected octant volume to be %v but got %v", expect, volume)
		}
	})
}

func TestGetChildAABCCases(t *testing.T) {
	testCases := []struct {
		desc   string
		curr   *chunk.AABC
		target glm.Vec3
		expect *chunk.AABC
	}{
		{
			desc: "GetChildAABC +x +y +z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{0, 0, 0},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, 0},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC +x +y -z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{0, 0, -1},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, 0, -2},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC +x -y +z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{0, -1, 0},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, -2, 0},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC +x -y -z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{0, -1, -1},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{0, -2, -2},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC -x +y +z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{-1, 0, 0},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{-2, 0, 0},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC -x +y -z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{-1, 0, -1},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{-2, 0, -2},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC -x -y +z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{-1, -1, 0},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, 0},
				Size: 2,
			},
		},
		{
			desc: "GetChildAABC -x -y -z",
			curr: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 4,
			},
			target: glm.Vec3{-1, -1, -1},
			expect: &chunk.AABC{
				Pos:  glm.Vec3{-2, -2, -2},
				Size: 2,
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			result := chunk.GetChildAABC(tC.curr, tC.target)
			if *result != *tC.expect {
				t.Fatalf("got AABC %v but expected %v", *result, *tC.expect)
			}
		})
	}
}
