package chunk_test

import (
	"testing"

	"github.com/kroppt/voxels/chunk"
)

func TestChunk(t *testing.T) {
	t.Parallel()
	t.Run("test get chunk position", func(t *testing.T) {
		t.Parallel()
		expected := chunk.Position{1, 2, 3}
		chunk := chunk.New(expected, 10)
		actual := chunk.Position()
		if actual != expected {
			t.Fatalf("expected to get chunk position %v but got %v", expected, actual)
		}
	})
	t.Run("test get chunk size", func(t *testing.T) {
		t.Parallel()
		expected := uint32(10)
		chunk := chunk.New(chunk.Position{1, 2, 3}, expected)
		actual := chunk.Size()
		if actual != expected {
			t.Fatalf("expected to get chunk size %v but got %v", expected, actual)
		}
	})
	t.Run("check that default block type is air", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeAir
		voxelCoordinate := chunk.VoxelCoordinate{4, 4, 4}
		chunk := chunk.New(chunk.Position{0, 0, 0}, 6)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("set block type of one voxel to dirt and get it back", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeDirt
		voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
		chunk := chunk.New(chunk.Position{0, 0, 0}, 10)
		chunk.SetBlockType(voxelCoordinate, expected)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
}

func TestChunkVoxelAdjacencies(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		adjMask chunk.AdjacentMask
		expect  chunk.AdjacentMask
	}{
		{
			desc:    "set front adj",
			adjMask: chunk.AdjacentFront,
			expect:  chunk.AdjacentFront,
		},
		{
			desc:    "set back adj",
			adjMask: chunk.AdjacentBack,
			expect:  chunk.AdjacentBack,
		},
		{
			desc:    "set left adj",
			adjMask: chunk.AdjacentLeft,
			expect:  chunk.AdjacentLeft,
		},
		{
			desc:    "set right adj",
			adjMask: chunk.AdjacentRight,
			expect:  chunk.AdjacentRight,
		},
		{
			desc:    "set top adj",
			adjMask: chunk.AdjacentTop,
			expect:  chunk.AdjacentTop,
		},
		{
			desc:    "set bottom adj",
			adjMask: chunk.AdjacentBottom,
			expect:  chunk.AdjacentBottom,
		},
		{
			desc:    "set X adj",
			adjMask: chunk.AdjacentX,
			expect:  chunk.AdjacentX,
		},
		{
			desc:    "set Y adj",
			adjMask: chunk.AdjacentY,
			expect:  chunk.AdjacentY,
		},
		{
			desc:    "set Z adj",
			adjMask: chunk.AdjacentZ,
			expect:  chunk.AdjacentZ,
		},
		{
			desc:    "set all adj",
			adjMask: chunk.AdjacentAll,
			expect:  chunk.AdjacentAll,
		},
		{
			desc:    "set none adj",
			adjMask: chunk.AdjacentNone,
			expect:  chunk.AdjacentNone,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
			chunk := chunk.New(chunk.Position{0, 0, 0}, 10)
			chunk.SetAdjacency(voxelCoordinate, tC.adjMask)
			actual := chunk.Adjacency(voxelCoordinate)
			if actual != tC.expect {
				t.Fatalf("expected to get adj mask %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestChunkVoxelLight(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		face   chunk.LightFace
		expect uint32
	}{
		{
			desc:   "set front light to 5",
			face:   chunk.LightFront,
			expect: 5,
		},
		{
			desc:   "set back light to 0",
			face:   chunk.LightFront,
			expect: 0,
		},
		{
			desc:   "set left light to 15",
			face:   chunk.LightLeft,
			expect: 15,
		},
		{
			desc:   "set right light to 8",
			face:   chunk.LightRight,
			expect: 8,
		},
		{
			desc:   "set bottom light to 6",
			face:   chunk.LightBottom,
			expect: 6,
		},
		{
			desc:   "set top light to 2",
			face:   chunk.LightTop,
			expect: 2,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
			chunk := chunk.New(chunk.Position{0, 0, 0}, 10)
			chunk.SetLighting(voxelCoordinate, tC.face, tC.expect)
			actual := chunk.Lighting(voxelCoordinate, tC.face)
			if actual != tC.expect {
				t.Fatalf("expected to get light value %v but got %v", tC.expect, actual)
			}
		})
	}
}
