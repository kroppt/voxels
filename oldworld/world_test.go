package oldworld_test

import (
	"testing"

	oldworld "github.com/kroppt/voxels/oldworld"
)

func TestGetChunkIndex(t *testing.T) {
	testCases := []struct {
		desc      string
		pos       oldworld.VoxelPos
		chunkSize int
		expect    oldworld.ChunkPos
	}{
		{
			desc:      "minimum edge inclusive",
			pos:       oldworld.VoxelPos{0, 0, 0},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{0, 0, 0},
		},
		{
			desc:      "maximum x edge exclusive",
			pos:       oldworld.VoxelPos{3, 2, 2},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{1, 0, 0},
		},
		{
			desc:      "maximum z edge exclusive",
			pos:       oldworld.VoxelPos{2, 2, 3},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{0, 0, 1},
		},
		{
			desc:      "negative first chunk x",
			pos:       oldworld.VoxelPos{-1, 1, 1},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{-1, 0, 0},
		},
		{
			desc:      "negative second chunk x",
			pos:       oldworld.VoxelPos{-5, 1, 1},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{-2, 0, 0},
		},
		{
			desc:      "negative y no change",
			pos:       oldworld.VoxelPos{1, -1, 1},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{0, -1, 0},
		},
		{
			desc:      "negative first chunk z",
			pos:       oldworld.VoxelPos{1, 1, -1},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{0, 0, -1},
		},
		{
			desc:      "negative second chunk z",
			pos:       oldworld.VoxelPos{1, 1, -5},
			chunkSize: 3,
			expect:    oldworld.ChunkPos{0, 0, -2},
		},
		{
			desc:      "regression",
			pos:       oldworld.VoxelPos{3, 2, 3},
			chunkSize: 6,
			expect:    oldworld.ChunkPos{0, 0, 0},
		},
		{
			desc:      "far negative edge is inclusive",
			pos:       oldworld.VoxelPos{0, 0, -5},
			chunkSize: 5,
			expect:    oldworld.ChunkPos{0, 0, -1},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actual := tC.pos.GetChunkPos(tC.chunkSize)
			if tC.expect != actual {
				t.Fatalf("expected %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestGetChunkBounds(t *testing.T) {
	testCases := []struct {
		desc       string
		renderDist int
		currChunk  oldworld.ChunkPos
		expectRng  oldworld.ChunkRange
	}{
		{
			desc:       "3x3 around 0,0",
			renderDist: 1,
			currChunk:  oldworld.ChunkPos{0, 0, 0},
			expectRng: oldworld.ChunkRange{
				Min: oldworld.ChunkPos{X: -1, Y: -1, Z: -1},
				Max: oldworld.ChunkPos{X: 1, Y: 1, Z: 1},
			},
		},
		{
			desc:       "5x5 around -1,-1",
			renderDist: 2,
			currChunk:  oldworld.ChunkPos{-1, 0, -1},
			expectRng: oldworld.ChunkRange{
				Min: oldworld.ChunkPos{X: -3, Y: -2, Z: -3},
				Max: oldworld.ChunkPos{X: 1, Y: 2, Z: 1},
			},
		},
		{
			desc:       "1x1 around 3,5",
			renderDist: 0,
			currChunk:  oldworld.ChunkPos{3, 0, 5},
			expectRng: oldworld.ChunkRange{
				Min: oldworld.ChunkPos{X: 3, Y: 0, Z: 5},
				Max: oldworld.ChunkPos{X: 3, Y: 0, Z: 5},
			},
		},
		{
			desc:       "3x3 around -1,2",
			renderDist: 1,
			currChunk:  oldworld.ChunkPos{-1, 0, 2},
			expectRng: oldworld.ChunkRange{
				Min: oldworld.ChunkPos{X: -2, Y: -1, Z: 1},
				Max: oldworld.ChunkPos{X: 0, Y: 1, Z: 3},
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actualRng := tC.currChunk.GetSurroundings(tC.renderDist)
			if tC.expectRng != actualRng {
				t.Fatalf("expected %v but got %v", tC.expectRng, actualRng)
			}
			if tC.expectRng != actualRng {
				t.Fatalf("expected %v but got %v", tC.expectRng, actualRng)
			}
		})
	}
}
