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
			actual := tC.pos.AsChunkPos(tC.chunkSize)
			if tC.expect != actual {
				t.Fatalf("expected %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestGetChunkBounds(t *testing.T) {
	testCases := []struct {
		desc       string
		renderDist int32
		currChunk  world.ChunkPos
		expectRng  world.ChunkRange
	}{
		{
			desc:       "3x3 around 0,0",
			renderDist: 1,
			currChunk:  world.ChunkPos{0, 0},
			expectRng: world.ChunkRange{
				Min: world.ChunkPos{X: -1, Z: -1},
				Max: world.ChunkPos{X: 1, Z: 1},
			},
		},
		{
			desc:       "5x5 around -1,-1",
			renderDist: 2,
			currChunk:  world.ChunkPos{-1, -1},
			expectRng: world.ChunkRange{
				Min: world.ChunkPos{X: -3, Z: -3},
				Max: world.ChunkPos{X: 1, Z: 1},
			},
		},
		{
			desc:       "1x1 around 3,5",
			renderDist: 0,
			currChunk:  world.ChunkPos{3, 5},
			expectRng: world.ChunkRange{
				Min: world.ChunkPos{X: 3, Z: 5},
				Max: world.ChunkPos{X: 3, Z: 5},
			},
		},
		{
			desc:       "3x3 around -1,2",
			renderDist: 1,
			currChunk:  world.ChunkPos{-1, 2},
			expectRng: world.ChunkRange{
				Min: world.ChunkPos{X: -2, Z: 1},
				Max: world.ChunkPos{X: 0, Z: 3},
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actualRng := world.GetChunkBounds(tC.renderDist, tC.currChunk)
			if tC.expectRng != actualRng {
				t.Fatalf("expected %v but got %v", tC.expectRng, actualRng)
			}
			if tC.expectRng != actualRng {
				t.Fatalf("expected %v but got %v", tC.expectRng, actualRng)
			}
		})
	}
}
