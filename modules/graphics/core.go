package graphics

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/util"
	"github.com/veandco/go-sdl2/sdl"
)

type core struct {
	window *sdl.Window
}

// ErrRenderDriver indicates that SDL failed to enable the OpenGL render driver.
const ErrRenderDriver log.ConstErr = "failed to set opengl render driver hint"

func (c *core) createWindow(title string, width, height uint32) error {
	runtime.LockOSThread()

	if !sdl.SetHint(sdl.HINT_RENDER_DRIVER, "opengl") {
		return ErrRenderDriver
	}
	var err error
	if err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		return err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4); err != nil {
		return err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3); err != nil {
		return err
	}
	if err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		return err
	}
	if err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE); err != nil {
		return err
	}

	var window *sdl.Window
	if window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height), sdl.WINDOW_HIDDEN|sdl.WINDOW_OPENGL); err != nil {
		return err
	}
	window.SetResizable(true)
	// creates context AND makes current
	if _, err = window.GLCreateContext(); err != nil {
		return err
	}
	if err = sdl.GLSetSwapInterval(1); err != nil {
		return err
	}

	if err = gl.Init(); err != nil {
		return err
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Debug("OpenGL version ", version)

	c.window = window
	return nil
}

func (c *core) showWindow() {
	c.window.Show()
}

func (c *core) pollEvent() (sdl.Event, bool) {
	evt := sdl.PollEvent()
	return evt, evt != nil
}

func (c *core) destroyWindow() error {
	err := c.window.Destroy()
	sdl.Quit()
	return err
}

func (c *core) updateView(viewableChunk map[chunk.Position]struct{}, viewMat mgl.Mat4) {

}

func (c *core) render() {
	w, h := c.window.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	sw := util.Start()
	// app.world.Render()
	sw.StopRecordAverage("Total World Render")
	c.window.GLSwap()

	for glErr := gl.GetError(); glErr != gl.NO_ERROR; glErr = gl.GetError() {
		log.Warnf("OpenGL error: %v", glErr)
	}
}
