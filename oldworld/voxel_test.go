package oldworld_test

import (
	"testing"

	oldworld "github.com/kroppt/voxels/oldworld"
)

func TestVbitSeparate(t *testing.T) {
	v := oldworld.Voxel{
		Pos:     oldworld.VoxelPos{},
		AdjMask: oldworld.AdjacentBack,
		Btype:   oldworld.Dirt,
	}
	vbits := v.GetVbits()
	adjMask, btype := oldworld.SeparateVbits(vbits)
	if adjMask != oldworld.AdjacentBack {
		t.Fatalf("expected adjmask %v but got %v", oldworld.AdjacentBack, adjMask)
	}
	if btype != oldworld.Dirt {
		t.Fatalf("expected block type %v but got %v", oldworld.Dirt, btype)
	}
}

func TestVbitConversion(t *testing.T) {
	v := oldworld.Voxel{
		Pos:     oldworld.VoxelPos{},
		AdjMask: oldworld.AdjacentBack,
		Btype:   oldworld.Dirt,
	}
	vbits := v.GetVbits()
	adjMask, btype := oldworld.SeparateVbits(vbits)
	v2 := oldworld.Voxel{
		Pos:     oldworld.VoxelPos{},
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
		mask      oldworld.LightMask
		expect    uint32
	}{
		{
			desc:      "set 0, get 0 left",
			lightInit: 0,
			setValue:  0,
			mask:      oldworld.LightLeft,
			expect:    0,
		},
		{
			desc:      "set 1, get 1 left",
			lightInit: 0,
			setValue:  1,
			mask:      oldworld.LightLeft,
			expect:    1,
		},
		{
			desc:      "set 2, get 2 back",
			lightInit: 0,
			setValue:  2,
			mask:      oldworld.LightBack,
			expect:    2,
		},
		{
			desc:      "set 3, get 3 front",
			lightInit: 0,
			setValue:  3,
			mask:      oldworld.LightFront,
			expect:    3,
		},
		{
			desc:      "set 4, get 4 bottom",
			lightInit: 0,
			setValue:  4,
			mask:      oldworld.LightBottom,
			expect:    4,
		},
		{
			desc:      "set 5, get 5 top",
			lightInit: 0,
			setValue:  5,
			mask:      oldworld.LightTop,
			expect:    5,
		},
		{
			desc:      "set 6, get 6 right",
			lightInit: 0,
			setValue:  6,
			mask:      oldworld.LightRight,
			expect:    6,
		},
		{
			desc:      "set 7, get 7 right",
			lightInit: 0,
			setValue:  7,
			mask:      oldworld.LightRight,
			expect:    7,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			v := oldworld.Voxel{
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
	v := oldworld.Voxel{}
	v.SetLightValue(3, oldworld.LightLeft)
	v.SetLightValue(7, oldworld.LightBottom)
	v.SetLightValue(5, oldworld.LightTop)
	if v.GetLightValue(oldworld.LightLeft) != 3 {
		t.Fatal("invalid light left")
	}
	if v.GetLightValue(oldworld.LightBottom) != 7 {
		t.Fatal("invalid light bottom")
	}
	if v.GetLightValue(oldworld.LightTop) != 5 {
		t.Fatal("invalid light top")
	}
}

func TestLightBitsOverwrite(t *testing.T) {
	v := oldworld.Voxel{}
	v.SetLightValue(3, oldworld.LightLeft)
	v.SetLightValue(7, oldworld.LightBottom)
	v.SetLightValue(5, oldworld.LightTop)
	v.SetLightValue(4, oldworld.LightLeft)

	if v.GetLightValue(oldworld.LightBottom) != 7 {
		t.Fatal("invalid light bottom")
	}
	if v.GetLightValue(oldworld.LightTop) != 5 {
		t.Fatal("invalid light top")
	}
	if v.GetLightValue(oldworld.LightLeft) != 4 {
		t.Fatal("invalid light left")
	}
}

func TestLbitSeparate(t *testing.T) {
	v := oldworld.Voxel{
		LightBits: uint32(oldworld.LightAll),
	}
	lbits := v.GetLbits()
	lightBits := oldworld.SeparateLbits(lbits)
	if lightBits != uint32(oldworld.LightAll) {
		t.Fatalf("expected lightBits %v but got %v", oldworld.LightAll, lightBits)
	}
}

func TestLbitConversion(t *testing.T) {
	v := oldworld.Voxel{
		LightBits: uint32(oldworld.LightAll),
	}
	lbits := v.GetLbits()
	lightBits := oldworld.SeparateLbits(lbits)
	v2 := oldworld.Voxel{
		LightBits: lightBits,
	}
	lbits2 := v2.GetLbits()
	if lbits != lbits2 {
		t.Fatalf("got lbits %v but wanted %v", lbits2, lbits)
	}
}
