package settings_test

import (
	"errors"
	"strings"
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
		expectX := uint32(1920)
		expectY := uint32(1080)

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

func TestRepositoryRenderDistance(t *testing.T) {
	t.Parallel()

	t.Run("set then get is same", func(t *testing.T) {
		t.Parallel()
		settings := settings.New()
		expected := uint32(10)

		settings.SetRenderDistance(expected)
		actual := settings.GetRenderDistance()

		if expected != actual {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})
}

func TestRepositoryFromReader(t *testing.T) {
	t.Parallel()

	t.Run("fails parsing equals signs", func(t *testing.T) {
		reader := strings.NewReader(strings.Join([]string{
			"resolutionX=100",
			"fov=60=60",
		}, "\n"))
		expect := settings.ErrParseSyntax
		var expectAs *settings.ErrParse
		expectLine := 2
		settingsMod := settings.New()

		err := settingsMod.SetFromReader(reader)

		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, expect) {
			t.Fatalf("expected %q but got %q", expect, err)
		}
		if !errors.As(err, &expectAs) {
			t.Fatalf("expected %T but got %T", expectAs, err)
		}
		if expectAs.Line != expectLine {
			t.Fatalf("expected %v but got %v", expectLine, expectAs.Line)
		}
	})

	t.Run("fails parsing invalid numbers", func(t *testing.T) {
		reader := strings.NewReader(strings.Join([]string{
			"fov=abc",
		}, "\n"))
		expect := settings.ErrParseValue
		var expectAs *settings.ErrParse
		expectLine := 1
		settingsMod := settings.New()

		err := settingsMod.SetFromReader(reader)

		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, expect) {
			t.Fatalf("expected %q but got %q", expect, err)
		}
		if !errors.As(err, &expectAs) {
			t.Fatalf("expected %T but got %T", expectAs, err)
		}
		if expectAs.Line != expectLine {
			t.Fatalf("expected %v but got %v", expectLine, expectAs.Line)
		}
	})

	t.Run("set and get is the same", func(t *testing.T) {
		reader := strings.NewReader(strings.Join([]string{
			"fov=60",
			"resolutionX=1920",
			"resolutionY=1080",
			"renderDistance=9",
		}, "\n"))
		settings := settings.New()

		err := settings.SetFromReader(reader)

		if err != nil {
			t.Fatal(err)
		}

		expectFOV := 60.0
		expectResolutionX := 1920
		expectResolutionY := 1080
		expectRenderDistance := 9

		fov := settings.GetFOV()
		if fov != expectFOV {
			t.Fatalf("expected %v but got %v", expectFOV, fov)
		}
		resolutionX, resolutionY := settings.GetResolution()
		if resolutionX != uint32(expectResolutionX) {
			t.Fatalf("expected %v but got %v", expectResolutionX, resolutionX)
		}
		if resolutionY != uint32(expectResolutionY) {
			t.Fatalf("expected %v but got %v", expectResolutionY, resolutionY)
		}
		renderDistance := settings.GetRenderDistance()
		if renderDistance != uint32(expectRenderDistance) {
			t.Fatalf("expected %v but got %v", expectRenderDistance, renderDistance)
		}
	})
}
