package app

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/shapes"
	"github.com/kroppt/voxels/voxgl"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win     *sdl.Window
	cube    *voxgl.Object
	running bool
}

func New(win *sdl.Window) (*Application, error) {
	colors := [8][3]float32{
		{0.0, 0.0, 1.0},
		{1.0, 0.0, 1.0},
		{0.0, 1.0, 1.0},
		{1.0, 1.0, 1.0},
		{0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{1.0, 1.0, 0.0},
	}
	cube, err := shapes.NewColoredCube(colors)
	if err != nil {
		return nil, err
	}

	return &Application{
		win:  win,
		cube: cube,
	}, nil
}

func (app *Application) Start() {
	app.running = true
	app.win.Show()
}

func (app *Application) Running() bool {
	return app.running
}

// HandleSdlEvent checks the type of a given SDL event and runs the method associated with that event
func (app *Application) HandleSdlEvent(e sdl.Event) {
	switch evt := e.(type) {
	case *sdl.QuitEvent:
		app.handleQuitEvent(evt)
	case *sdl.MouseButtonEvent:
		//app.handleMouseButtonEvent(evt)
	case *sdl.MouseMotionEvent:
		//app.handleMouseMotionEvent(evt)
	case *sdl.MouseWheelEvent:
		//app.handleMouseWheelEvent(evt)
	case *sdl.WindowEvent:
		//app.handleWindowEvent(evt)
	case *sdl.SysWMEvent:
		//app.handleSysWMEvent(evt)
	case *sdl.KeyboardEvent:
		app.handleKeyboardEvent(evt)
	}
}

func (app *Application) handleQuitEvent(evt *sdl.QuitEvent) {
	app.running = false
}

func (app *Application) handleKeyboardEvent(evt *sdl.KeyboardEvent) {
	if evt.State != sdl.PRESSED {
		return
	}
	deg := float32(2.0)
	switch evt.Keysym.Sym {
	case sdl.K_UP:
		app.cube.Rotate(-deg, 0, 0)
	case sdl.K_DOWN:
		app.cube.Rotate(deg, 0, 0)
	case sdl.K_LEFT:
		app.cube.Rotate(0, -deg, 0)
	case sdl.K_RIGHT:
		app.cube.Rotate(0, deg, 0)
	}
}

func (app *Application) PostEventActions() {
	w, h := app.win.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	app.cube.Render()

	app.win.GLSwap()

	for glErr := gl.GetError(); glErr != gl.NO_ERROR; glErr = gl.GetError() {
		log.Warnf("OpenGL error: %v", glErr)
	}
}

func (app *Application) Quit() {
	if err := app.win.Destroy(); err != nil {
		log.Fatal(err)
	}
	app.cube.Destroy()
	sdl.Quit()
}
