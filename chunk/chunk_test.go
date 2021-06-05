package chunk_test

import (
	"testing"

	"github.com/engoengine/glm"
	"github.com/kroppt/voxels/chunk"
)

func TestNewChunk(t *testing.T) {
	chunk.NewChunk(0, 0, glm.Vec2{})
}
