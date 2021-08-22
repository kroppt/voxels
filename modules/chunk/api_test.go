package chunk_test

import (
	"testing"

	"github.com/kroppt/voxels/modules/chunk"
)

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := chunk.New()
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}
