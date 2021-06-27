package world_test

import (
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestNewChunkNotNil(t *testing.T) {
	ch := world.NewChunk2()
	if ch == nil {
		t.Fatalf("new chunk was nil")
	}
}
