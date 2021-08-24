package settings_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kroppt/voxels/repositories/settings"
)

type fnFileMod struct {
	fnGetFileReader func(fileName string) io.Reader
}

func (fn *fnFileMod) GetFileReader(fileName string) io.Reader {
	return fn.fnGetFileReader(fileName)
}

func TestRepositoryNew(t *testing.T) {
	t.Parallel()

	t.Run("should not return nil", func(t *testing.T) {
		t.Parallel()

		settings := settings.New(nil)

		if settings == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestRepositoryFOV(t *testing.T) {
	t.Parallel()

	t.Run("set then get is same", func(t *testing.T) {
		t.Parallel()

		settings := settings.New(nil)
		expect := float32(90)

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

		settings := settings.New(nil)
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

func TestRepositoryFromFile(t *testing.T) {
	t.Parallel()

	t.Run("fails correctly", func(t *testing.T) {
		fileMod := &fnFileMod{
			fnGetFileReader: func(fileName string) io.Reader {
				if fileName != "settings.config" {
					return nil
				}
				return strings.NewReader(strings.Join([]string{
					"fov=60=60",
				}, "\n"))
			},
		}
		settingsMod := settings.New(fileMod)

		reader := fileMod.GetFileReader("settings.config")

		err := settingsMod.SetFromReader(reader)

		if !errors.Is(err, settings.ErrSettingsParse) {
			t.Fatalf("expected %v but got %v", settings.ErrSettingsParse, err)
		}
	})

	t.Run("set and get is the same", func(t *testing.T) {
		fileMod := &fnFileMod{
			fnGetFileReader: func(fileName string) io.Reader {
				if fileName != "settings.config" {
					return nil
				}
				return strings.NewReader(strings.Join([]string{
					"fov=60",
					"resolutionX=1920",
					"resolutionY=1080",
				}, "\n"))
			},
		}
		settings := settings.New(fileMod)

		reader := fileMod.GetFileReader("settings.config")

		err := settings.SetFromReader(reader)

		if err != nil {
			t.Fatal(err)
		}

		expectFOV := 60
		expectResolutionX := 1920
		expectResolutionY := 1080

		fov := settings.GetFOV()
		if fov != float32(expectFOV) {
			t.Fatalf("expected %v but got %v", expectFOV, fov)
		}
		resolutionX, resolutionY := settings.GetResolution()
		if resolutionX != int32(expectResolutionX) {
			t.Fatalf("expected %v but got %v", expectResolutionX, resolutionX)
		}
		if resolutionY != int32(expectResolutionY) {
			t.Fatalf("expected %v but got %v", expectResolutionY, resolutionY)
		}
	})
}
