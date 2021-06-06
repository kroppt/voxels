package world_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/world"
)

func TestNewChunk(t *testing.T) {
	world.NewChunk(0, 0, glm.Vec2{})
}

func TestIsWithinChunk(t *testing.T) {
	t.Run("standard within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{1, 0, -7}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("minimum within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{0, 0, -10}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("maximum within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{4, 0, -6}
		result := chunk.IsWithinChunk(pos)
		expect := true
		if result != expect {
			t.Fatalf("Expected %v to be in chunk, but was not", pos)
		}
	})
	t.Run("maximum z out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{4, 0, -5}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("maximum x out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{5, 0, -6}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("negative y out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{2, -1, -8}
		result := chunk.IsWithinChunk(pos)
		expect := false
		if result != expect {
			t.Fatalf("Expected %v to out be in chunk, but it was inside", pos)
		}
	})
	t.Run("too large y out of bounds within chunk test", func(t *testing.T) {
		t.Parallel()
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{2, 1, -8}
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
		chunk := world.NewChunk(5, 1, glm.Vec2{0, -2})
		pos := world.VoxelPos{1, 3, -7}
		i, j, k := chunk.GetRelativeIndices(pos)
		if i != 1 || j != 3 || k != 3 {
			t.Fatalf("expected 1, 3, 3 but got %v %v %v", i, j, k)
		}
	})
}
