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
