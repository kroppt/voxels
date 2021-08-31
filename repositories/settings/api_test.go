package settings_test

import (
	"testing"

	"github.com/kroppt/voxels/repositories/settings"
)

func TestRepositoryNew(t *testing.T) {
	t.Parallel()

	t.Run("should not return nil", func(t *testing.T) {
		t.Parallel()

		settings := settings.New()

		if settings == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestRepositoryFOV(t *testing.T) {
	t.Parallel()

	t.Run("set then get is same", func(t *testing.T) {
		t.Parallel()

		settings := settings.New()
		expect := 90.0

		settings.SetFOV(expect)
		got := settings.GetFOV()

		if got != expect {
			t.Fatalf("expected %v but got %v", expect, got)
		}
	})
}

func TestRepositoryResolution(t *testing.T) {
	t.Parallel()

	t.Run("set then get is same", func(t *testing.T) {
		t.Parallel()

		settings := settings.New()
		expectX := int32(1920)
		expectY := int32(1080)

		settings.SetResolution(expectX, expectY)
		gotX, gotY := settings.GetResolution()

		if gotX != expectX {
			t.Fatalf("expected %v but got %v", expectX, gotX)
		}
		if gotY != expectY {
			t.Fatalf("expected %v but got %v", expectY, gotY)
		}
	})
}
