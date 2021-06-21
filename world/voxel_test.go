package world_test

import (
	"testing"

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

func TestLightBitsOverwriteWithValue(t *testing.T) {
	v := world.Voxel{}
	v.SetLightValue(3, world.LightLeft)
	v.SetLightValue(7, world.LightBottom)
	v.SetLightValue(5, world.LightTop)
	v.SetLightValue(4, world.LightLeft)
	v.SetLightValue(2, world.LightValue)

	if v.GetLightValue(world.LightBottom) != 7 {
		t.Fatal("invalid light bottom")
	}
	if v.GetLightValue(world.LightTop) != 5 {
		t.Fatal("invalid light top")
	}
	if v.GetLightValue(world.LightLeft) != 4 {
		t.Fatal("invalid light left")
	}
	v.SetLightValue(0, world.LightValue)

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
