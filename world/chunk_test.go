package world_test

import (
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestNewChunk(t *testing.T) {
	world.NewChunk(0, world.ChunkPos{}, world.FlatWorldGenerator{})
}

func TestIsWithinChunk(t *testing.T) {
	t.Run("standard within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{1, 0, -7}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("minimum within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{0, 0, -10}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("maximum within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{4, 0, -6}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("maximum z out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{4, 0, -5}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("maximum x out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{5, 0, -6}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("negative y out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{2, -1, -8}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("too large y out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{2, 5, -8}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
}

func TestGetRelativeIndices(t *testing.T) {
	t.Run("", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, world.ChunkPos{0, 0, -2}, world.FlatWorldGenerator{})
		pos := world.VoxelPos{1, 3, -7}
		localPos := pos.AsLocalChunkPos(chunk)
		i, j, k := localPos.X, localPos.Y, localPos.Z
		if i != 1 || j != 3 || k != 3 {
			t.Fatalf("expected 1, 3, 3 but got %v %v %v", i, j, k)
		}
	})
}

func TestChunkRange(t *testing.T) {
	results := make(map[world.ChunkPos]struct{})
	rng := world.ChunkRange{
		Min: world.ChunkPos{-3, -1, 0},
		Max: world.ChunkPos{1, 1, 6},
	}
	rng.ForEach(func(pos world.ChunkPos) bool {
		results[pos] = struct{}{}
		return false
	})
	for i := -3; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := 0; k <= 6; k++ {
				if _, ok := results[world.ChunkPos{i, j, k}]; !ok {
					t.Fatalf("(%v, %v, %v) was expected to in map, but wasn't", i, j, k)
				}
			}
		}
	}
}

func TestChunkRangeSub(t *testing.T) {
	testCases := []struct {
		desc    string
		new     world.ChunkRange
		old     world.ChunkRange
		nExpect int
	}{
		{
			desc: "8/8 ranges",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{4, 4, 4},
				Max: world.ChunkPos{6, 6, 6},
			},
			nExpect: 8,
		},
		{
			desc: "7/8 ranges",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{1, 1, 1},
				Max: world.ChunkPos{3, 3, 3},
			},
			nExpect: 7,
		},
		{
			desc: "7/8 ranges other diagonal",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{-1, -1, -1},
				Max: world.ChunkPos{1, 1, 1},
			},
			nExpect: 7,
		},
		{
			desc: "7/8 ranges other diagonal #2",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{1, 1, -1},
				Max: world.ChunkPos{3, 3, 1},
			},
			nExpect: 7,
		},
		{
			desc: "6/8 ranges",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{1, 0, 1},
				Max: world.ChunkPos{3, 2, 3},
			},
			nExpect: 6,
		},
		{
			desc: "6/8 ranges another way",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{-1, 0, -1},
				Max: world.ChunkPos{1, 2, 1},
			},
			nExpect: 6,
		},
		{
			desc: "4/8 ranges",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 1},
				Max: world.ChunkPos{2, 2, 3},
			},
			nExpect: 4,
		},
		{
			desc: "0/8 ranges",
			new: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			old: world.ChunkRange{
				Min: world.ChunkPos{0, 0, 0},
				Max: world.ChunkPos{2, 2, 2},
			},
			nExpect: 0,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			ranges := tC.new.Sub(tC.old)
			nRanges := len(ranges)
			if nRanges != tC.nExpect {
				t.Fatalf("expected %v ranges but got %v", tC.nExpect, nRanges)
			}
		})
	}
}

func TestForEachSub7(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{1, 1, 1},
		Max: world.ChunkPos{3, 3, 3},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{
		{1, 1, 2},
		{1, 2, 1},
		{1, 2, 2},
		{2, 1, 1},
		{2, 1, 2},
		{2, 2, 1},
		{2, 2, 2},
	}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestForEachSub6(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{1, 0, 1},
		Max: world.ChunkPos{3, 2, 3},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{
		{1, 0, 2},
		{1, 1, 2},
		{2, 0, 1},
		{2, 1, 2},
		{2, 0, 2},
		{2, 1, 2},
	}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestForEachSub4(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 1},
		Max: world.ChunkPos{2, 2, 3},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{
		{0, 0, 2},
		{0, 1, 2},
		{1, 0, 2},
		{1, 1, 2},
	}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestForEachSubNone(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestForEachSub8(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{2, 2, 2},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{-2, -2, -2},
		Max: world.ChunkPos{0, 0, 0},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{
		{-2, -2, -2},
		{-2, -2, -1},
		{-2, -1, -2},
		{-2, -1, -1},
		{-1, -2, -2},
		{-1, -2, -1},
		{-1, -1, -2},
		{-1, -1, -1},
	}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestForEachSub8EmptyOld(t *testing.T) {
	old := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{0, 0, 0},
	}
	new := world.ChunkRange{
		Min: world.ChunkPos{-2, -2, -2},
		Max: world.ChunkPos{0, 0, 0},
	}
	result := make(map[world.ChunkPos]struct{})
	expectKeys := []world.ChunkPos{
		{-2, -2, -2},
		{-2, -2, -1},
		{-2, -1, -2},
		{-2, -1, -1},
		{-1, -2, -2},
		{-1, -2, -1},
		{-1, -1, -2},
		{-1, -1, -1},
	}
	new.ForEachSub(old, func(pos world.ChunkPos) bool {
		result[pos] = struct{}{}
		return false
	})
	if len(result) != len(expectKeys) {
		t.Fatalf("expected %v keys but found %v", len(expectKeys), len(result))
	}
	for _, key := range expectKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected chunk pos %v to be in map, but was not", key)
		}
	}
}

func TestChunkRangeCount2x2x2(t *testing.T) {
	rng := world.ChunkRange{
		Min: world.ChunkPos{-2, -2, -2},
		Max: world.ChunkPos{0, 0, 0},
	}
	expected := 8
	result := rng.Count()
	if result != expected {
		t.Fatalf("expected %v but got %v", expected, result)
	}
}

func TestChunkRangeCount1x1x1(t *testing.T) {
	rng := world.ChunkRange{
		Min: world.ChunkPos{-1, -1, -1},
		Max: world.ChunkPos{0, 0, 0},
	}
	expected := 1
	result := rng.Count()
	if result != expected {
		t.Fatalf("expected %v but got %v", expected, result)
	}
}

func TestChunkRangeCount0x0x0(t *testing.T) {
	rng := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{0, 0, 0},
	}
	expected := 0
	result := rng.Count()
	if result != expected {
		t.Fatalf("expected %v but got %v", expected, result)
	}
}

func TestChunkRangeCount1x2x3(t *testing.T) {
	rng := world.ChunkRange{
		Min: world.ChunkPos{0, 0, 0},
		Max: world.ChunkPos{1, 2, 3},
	}
	expected := 6
	result := rng.Count()
	if result != expected {
		t.Fatalf("expected %v but got %v", expected, result)
	}
}
