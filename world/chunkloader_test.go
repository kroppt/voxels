package world_test

import (
	"testing"

	"github.com/kroppt/voxels/world"
)

func TestNewChunkLoader(t *testing.T) {
	cl := world.NewChunkLoader()
	if cl == nil {
		t.Fatalf("new chunk loader was nil")
	}
}
