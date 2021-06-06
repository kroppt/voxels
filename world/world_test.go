package world_test

import (
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestGetChunkIndex(t *testing.T) {
	testCases := []struct {
		desc      string
		pos       world.VoxelPos
		chunkSize int32
		expect    world.ChunkPos
	}{
		{
			desc:      "minimum edge inclusive",
			pos:       world.VoxelPos{0, 0, 0},
			chunkSize: 3,
			expect:    world.ChunkPos{0, 0},
		},
		{
			desc:      "maximum x edge exclusive",
			pos:       world.VoxelPos{3, 2, 2},
			chunkSize: 3,
			expect:    world.ChunkPos{1, 0},
		},
		{
			desc:      "maximum z edge exclusive",
			pos:       world.VoxelPos{2, 2, 3},
			chunkSize: 3,
			expect:    world.ChunkPos{0, 1},
		},
		{
			desc:      "negative first chunk x",
			pos:       world.VoxelPos{-1, 1, 1},
			chunkSize: 3,
			expect:    world.ChunkPos{-1, 0},
		},
		{
			desc:      "negative second chunk x",
			pos:       world.VoxelPos{-5, 1, 1},
			chunkSize: 3,
			expect:    world.ChunkPos{-2, 0},
		},
		{
			desc:      "negative y no change",
			pos:       world.VoxelPos{1, -1, 1},
			chunkSize: 3,
			expect:    world.ChunkPos{0, 0},
		},
		{
			desc:      "negative first chunk z",
			pos:       world.VoxelPos{1, 1, -1},
			chunkSize: 3,
			expect:    world.ChunkPos{0, -1},
		},
		{
			desc:      "negative second chunk z",
			pos:       world.VoxelPos{1, 1, -5},
			chunkSize: 3,
			expect:    world.ChunkPos{0, -2},
		},
		{
			desc:      "regression",
			pos:       world.VoxelPos{3, 2, 3},
			chunkSize: 6,
			expect:    world.ChunkPos{0, 0},
		},
		{
			desc:      "far negative edge is inclusive",
			pos:       world.VoxelPos{0, 0, -5},
			chunkSize: 5,
			expect:    world.ChunkPos{0, -1},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actual := world.GetChunkIndex(tC.chunkSize, tC.pos)
			if tC.expect != actual {
				t.Fatalf("expected %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestGetChunkBounds(t *testing.T) {
	testCases := []struct {
		desc      string
		worldSize int32
		currChunk world.ChunkPos
		expectX   world.Range
		expectZ   world.Range
	}{
		{
			desc:      "3x3 around 0,0",
			worldSize: 3,
			currChunk: world.ChunkPos{0, 0},
			expectX:   world.Range{-1, 1},
			expectZ:   world.Range{-1, 1},
		},
		{
			desc:      "5x5 around -1,-1",
			worldSize: 5,
			currChunk: world.ChunkPos{-1, -1},
			expectX:   world.Range{-3, 1},
			expectZ:   world.Range{-3, 1},
		},
		{
			desc:      "1x1 around 3,5",
			worldSize: 1,
			currChunk: world.ChunkPos{3, 5},
			expectX:   world.Range{3, 3},
			expectZ:   world.Range{5, 5},
		},
		{
			desc:      "3x3 around -1,2",
			worldSize: 3,
			currChunk: world.ChunkPos{-1, 2},
			expectX:   world.Range{-2, 0},
			expectZ:   world.Range{1, 3},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actualX, actualZ := world.GetChunkBounds(tC.worldSize, tC.currChunk)
			if tC.expectX != actualX {
				t.Fatalf("expected %v but got %v", tC.expectX, actualX)
			}
			if tC.expectZ != actualZ {
				t.Fatalf("expected %v but got %v", tC.expectZ, actualZ)
			}
		})
	}
}
