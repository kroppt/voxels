package world_test

import (
	"fmt"
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/world"
)

func TestVbitSeparate(t *testing.T) {
	v := world.Voxel{
		Pos:     world.VoxelPos{},
		AdjMask: world.AdjacentBack,
		Btype:   world.Dirt,
	}
	vbits := v.GetVbits()
	adjMask, btype := world.SeparateVbits(vbits)
	if adjMask != world.AdjacentBack {
		t.Fatalf("expected adjmask %v but got %v", world.AdjacentBack, adjMask)
	}
	if btype != world.Dirt {
		t.Fatalf("expected block type %v but got %v", world.Dirt, btype)
	}
}

func TestVbitConversion(t *testing.T) {
	v := world.Voxel{
		Pos:     world.VoxelPos{},
		AdjMask: world.AdjacentBack,
		Btype:   world.Dirt,
	}
	vbits := v.GetVbits()
	adjMask, btype := world.SeparateVbits(vbits)
	v2 := world.Voxel{
		Pos:     world.VoxelPos{},
		AdjMask: adjMask,
		Btype:   btype,
	}
	vbits2 := v2.GetVbits()
	if vbits != vbits2 {
		t.Fatalf("got vbits %v but wanted %v", vbits2, vbits)
	}
}

func TestLightBits(t *testing.T) {
	testCases := []struct {
		desc      string
		lightInit uint32
		setValue  uint32
		mask      world.LightMask
		expect    uint32
	}{
		{
			desc:      "set 0, get 0 left",
			lightInit: 0,
			setValue:  0,
			mask:      world.LightLeft,
			expect:    0,
		},
		{
			desc:      "set 1, get 1 left",
			lightInit: 0,
			setValue:  1,
			mask:      world.LightLeft,
			expect:    1,
		},
		{
			desc:      "set 2, get 2 back",
			lightInit: 0,
			setValue:  2,
			mask:      world.LightBack,
			expect:    2,
		},
		{
			desc:      "set 3, get 3 front",
			lightInit: 0,
			setValue:  3,
			mask:      world.LightFront,
			expect:    3,
		},
		{
			desc:      "set 4, get 4 bottom",
			lightInit: 0,
			setValue:  4,
			mask:      world.LightBottom,
			expect:    4,
		},
		{
			desc:      "set 5, get 5 top",
			lightInit: 0,
			setValue:  5,
			mask:      world.LightTop,
			expect:    5,
		},
		{
			desc:      "set 6, get 6 right",
			lightInit: 0,
			setValue:  6,
			mask:      world.LightRight,
			expect:    6,
		},
		{
			desc:      "set 7, get 7 right",
			lightInit: 0,
			setValue:  7,
			mask:      world.LightRight,
			expect:    7,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC := tC
			// t.Parallel()
			v := world.Voxel{
				LightBits: tC.lightInit,
			}
			v.SetLightValue(tC.setValue, tC.mask)
			result := v.GetLightValue(tC.mask)
			if result != tC.expect {
				t.Fatalf("expected light value %v for mask %v, but got %v", tC.expect, tC.mask, result)
			}
		})
	}
}

func TestLightBitsCorruption(t *testing.T) {
	v := world.Voxel{}
	v.SetLightValue(3, world.LightLeft)
	v.SetLightValue(7, world.LightBottom)
	v.SetLightValue(5, world.LightTop)
	if v.GetLightValue(world.LightLeft) != 3 {
		t.Fatal("invalid light left")
	}
	if v.GetLightValue(world.LightBottom) != 7 {
		t.Fatal("invalid light bottom")
	}
	if v.GetLightValue(world.LightTop) != 5 {
		t.Fatal("invalid light top")
	}
}

func TestLightBitsOverwrite(t *testing.T) {
	v := world.Voxel{}
	v.SetLightValue(3, world.LightLeft)
	v.SetLightValue(7, world.LightBottom)
	v.SetLightValue(5, world.LightTop)
	v.SetLightValue(4, world.LightLeft)

	if v.GetLightValue(world.LightBottom) != 7 {
		t.Fatal("invalid light bottom")
	}
	if v.GetLightValue(world.LightTop) != 5 {
		t.Fatal("invalid light top")
	}
	if v.GetLightValue(world.LightLeft) != 4 {
		t.Fatal("invalid light left")
	}
}

func TestLbitSeparate(t *testing.T) {
	v := world.Voxel{
		LightBits: uint32(world.LightAll),
	}
	lbits := v.GetLbits()
	lightBits := world.SeparateLbits(lbits)
	if lightBits != uint32(world.LightAll) {
		t.Fatalf("expected lightBits %v but got %v", world.LightAll, lightBits)
	}
}

func TestLbitConversion(t *testing.T) {
	v := world.Voxel{
		LightBits: uint32(world.LightAll),
	}
	lbits := v.GetLbits()
	lightBits := world.SeparateLbits(lbits)
	v2 := world.Voxel{
		LightBits: lightBits,
	}
	lbits2 := v2.GetLbits()
	if lbits != lbits2 {
		t.Fatalf("got lbits %v but wanted %v", lbits2, lbits)
	}
}

func TestGetChunkPos(t *testing.T) {
	testCases := []struct {
		vpos      world.VoxelPos
		chunkSize int
		expected  world.ChunkPos
	}{
		{
			vpos:      world.VoxelPos{0, 0, 0},
			chunkSize: 1,
			expected:  world.ChunkPos{0, 0, 0},
		},
		{
			vpos:      world.VoxelPos{1, 0, 0},
			chunkSize: 1,
			expected:  world.ChunkPos{1, 0, 0},
		},
		{
			vpos:      world.VoxelPos{1, 0, 0},
			chunkSize: 2,
			expected:  world.ChunkPos{0, 0, 0},
		},
		{
			vpos:      world.VoxelPos{-1, 0, 0},
			chunkSize: 5,
			expected:  world.ChunkPos{-1, 0, 0},
		},
		{
			vpos:      world.VoxelPos{5, 5, 5},
			chunkSize: 3,
			expected:  world.ChunkPos{1, 1, 1},
		},
	}
	for _, tC := range testCases {
		tC := tC
		desc := fmt.Sprintf("voxel: %v, chunkSize: %v", tC.vpos, tC.chunkSize)
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			result := tC.vpos.GetChunkPos(tC.chunkSize)
			if result != tC.expected {
				t.Fatalf("expected %v but got %v", tC.expected, result)
			}
		})
	}
}

func TestVoxelArithmetic(t *testing.T) {
	t.Run("voxel add", func(t *testing.T) {
		t.Parallel()
		vpos := world.VoxelPos{-2, 0, 2}
		diff := world.VoxelPos{10, -5, 0}
		result := vpos.Add(diff)
		expected := world.VoxelPos{8, -5, 2}
		if result != expected {
			t.Fatalf("expected %v but got %v", expected, result)
		}
	})
	t.Run("voxel sub", func(t *testing.T) {
		t.Parallel()
		vpos := world.VoxelPos{-2, 0, 2}
		diff := world.VoxelPos{10, -5, 0}
		result := vpos.Sub(diff)
		expected := world.VoxelPos{-12, 5, 2}
		if result != expected {
			t.Fatalf("expected %v but got %v", expected, result)
		}
	})
	t.Run("voxel to glm.vec3", func(t *testing.T) {
		t.Parallel()
		vpos := world.VoxelPos{-2, 5, 2}
		result := vpos.AsVec3()
		expected := glm.Vec3{-2.0, 5.0, 2.0}
		if result != expected {
			t.Fatalf("expected %v but got %v", expected, result)
		}
	})
	t.Run("voxel ForEach sum", func(t *testing.T) {
		t.Parallel()
		sum := 15
		vpos := world.VoxelPos{0, 0, 0}
		rng := vpos.GetSurroundings(2)
		rng.ForEach(func(pos world.VoxelPos) {
			sum += pos.X + pos.Y + pos.Z
		})
		expected := 15
		if sum != expected {
			t.Fatalf("expected %v but got %v", expected, sum)
		}
	})
}
