package app

import (
	"runtime"
	"testing"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func BenchmarkAppFindLookatVoxel(b *testing.B) {
	runtime.LockOSThread()
	win, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		0, 0, sdl.WINDOW_HIDDEN|sdl.WINDOW_OPENGL)
	if err != nil {
		b.Fatal(err)
	}
	_, err = win.GLCreateContext()
	if err != nil {
		b.Fatal(err)
	}
	err = gl.Init()
	if err != nil {
		b.Fatal(err)
	}
	app, err := New(win)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		app.findLookatVoxel()
	}
}
