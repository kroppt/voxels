package settings_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kroppt/voxels/repositories/settings"
)

type fnFileMod struct {
	fnGetReadCloser func(fileName string) (io.ReadCloser, error)
}

func (fn *fnFileMod) GetReadCloser(fileName string) (io.ReadCloser, error) {
	return fn.fnGetReadCloser(fileName)
}

type fnReadCloser struct {
	fnRead  func(p []byte) (n int, err error)
	fnClose func() error
}

func (fn *fnReadCloser) Read(p []byte) (n int, err error) {
	return fn.fnRead(p)
}

func (fn *fnReadCloser) Close() error {
	return fn.fnClose()
}

const invalidFileNameError string = "invalid file name"

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

func TestRepositoryFromReader(t *testing.T) {
	t.Parallel()

	t.Run("fails parsing equals signs", func(t *testing.T) {
		fileMod := &fnFileMod{
			fnGetReadCloser: func(fileName string) (io.ReadCloser, error) {
				if fileName != "settings.config" {
					return nil, errors.New(invalidFileNameError)
				}
				reader := strings.NewReader(strings.Join([]string{
					"resolutionX=100",
					"fov=60=60",
				}, "\n"))
				return &fnReadCloser{
					fnRead: func(p []byte) (int, error) {
						return reader.Read(p)
					},
					fnClose: func() error {
						return nil
					},
				}, nil
			},
		}
		expect := settings.ErrParseSyntax
		var expectAs *settings.ErrParse
		expectLine := 2
		settingsMod := settings.New()
		readerCloser, err := fileMod.GetReadCloser("settings.config")
		if err != nil {
			t.Fatal(err)
		}

		err = settingsMod.SetFromReader(readerCloser)

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
		fileMod := &fnFileMod{
			fnGetReadCloser: func(fileName string) (io.ReadCloser, error) {
				if fileName != "settings.config" {
					return nil, errors.New(invalidFileNameError)
				}
				reader := strings.NewReader(strings.Join([]string{
					"fov=abc",
				}, "\n"))
				return &fnReadCloser{
					fnRead: func(p []byte) (int, error) {
						return reader.Read(p)
					},
					fnClose: func() error {
						return nil
					},
				}, nil
			},
		}
		expect := settings.ErrParseValue
		var expectAs *settings.ErrParse
		expectLine := 1
		settingsMod := settings.New()
		readCloser, err := fileMod.GetReadCloser("settings.config")
		if err != nil {
			t.Fatal(err)
		}

		err = settingsMod.SetFromReader(readCloser)

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
		fileMod := &fnFileMod{
			fnGetReadCloser: func(fileName string) (io.ReadCloser, error) {
				if fileName != "settings.config" {
					return nil, errors.New(invalidFileNameError)
				}
				reader := strings.NewReader(strings.Join([]string{
					"fov=60",
					"resolutionX=1920",
					"resolutionY=1080",
				}, "\n"))
				return &fnReadCloser{
					fnRead: func(p []byte) (int, error) {
						return reader.Read(p)
					},
					fnClose: func() error {
						return nil
					},
				}, nil
			},
		}
		settings := settings.New()

		readCloser, err := fileMod.GetReadCloser("settings.config")
		if err != nil {
			t.Fatal(err)
		}

		err = settings.SetFromReader(readCloser)

		if err != nil {
			t.Fatal(err)
		}

		expectFOV := 60.0
		expectResolutionX := 1920
		expectResolutionY := 1080

		fov := settings.GetFOV()
		if fov != expectFOV {
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
