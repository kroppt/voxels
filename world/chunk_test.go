package world_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/world"
)

func TestNewChunk(t *testing.T) {
	world.NewChunk(0, 0, glm.Vec2{})
}

func TestGetRelativeIndices(t *testing.T) {
	t.Run("", func(t *testing.T) {
		t.Parallel()
		pos := glm.Vec3{1, 3, -7}
		chunkPos := glm.Vec2{0, -10}
		i, j, k := world.GetRelativeIndices(chunkPos, pos)
		if i != 1 || j != 3 || k != 3 {
			t.Fatalf("expected 1, 3, 3 but got %v %v %v", i, j, k)
		}
	})
}
