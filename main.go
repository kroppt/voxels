package main

import (
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/appchunk"
	"github.com/kroppt/voxels/log"
	"github.com/veandco/go-sdl2/sdl"
)

// ErrRenderDriver indicates that SDL failed to enable the OpenGL render driver
const ErrRenderDriver log.ConstErr = "failed to set opengl render driver hint"

func initWindow(title string, width, height int32) (*sdl.Window, error) {
	if !sdl.SetHint(sdl.HINT_RENDER_DRIVER, "opengl") {
		return nil, ErrRenderDriver
	}
	var err error
	if err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		return nil, err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4); err != nil {
		return nil, err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3); err != nil {
		return nil, err
	}
	if err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		return nil, err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE); err != nil {
		return nil, err
	}

	var window *sdl.Window
	if window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		width, height, sdl.WINDOW_HIDDEN|sdl.WINDOW_OPENGL); err != nil {
		return nil, err
	}
	window.SetResizable(true)
	// creates context AND makes current
	if _, err = window.GLCreateContext(); err != nil {
		return nil, err
	}
	if err = sdl.GLSetSwapInterval(1); err != nil {
		return nil, err
	}

	if err = gl.Init(); err != nil {
		return nil, err
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Debug("OpenGL version ", version)

	return window, nil
}

func init() {
	runtime.LockOSThread()
}

func main() {
	log.SetInfoOutput(os.Stderr)
	log.SetWarnOutput(os.Stderr)
	log.SetDebugOutput(os.Stderr)
	log.SetPerfOutput(os.Stderr)
	log.SetFatalOutput(os.Stderr)
	log.SetColorized(true)

	win, err := initWindow("voxels", 1920, 1080)
	if err != nil {
		log.Fatal(err)
	}

	app, err := appchunk.New(win)
	if err != nil {
		log.Fatal(err)
	}
	app.Start()

	for app.Running() {
		for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			err := app.HandleSdlEvent(evt)
			if err != nil {
				log.Fatalf("error handling sdl events: %v", err)
			}
		}

		app.PostEventActions()
	}

	app.Quit()
}
