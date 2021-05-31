package app

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/world"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win           *sdl.Window
	plane         *world.Plane
	planeRenderer *PlaneRenderer
	running       bool
}

func New(win *sdl.Window) (*Application, error) {
	planeRenderer := NewPlaneRenderer()

	// 11 x 11 x 11
	x := world.Range{Min: -5, Max: 5}
	y := world.Range{Min: -5, Max: 5}
	z := world.Range{Min: -5, Max: 5}
	plane, err := world.NewPlane(planeRenderer, x, y, z)
	if err != nil {
		return nil, fmt.Errorf("could not create plane: %v", err)
	}

	return &Application{
		win:           win,
		plane:         plane,
		planeRenderer: planeRenderer,
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
	cam := app.plane.GetCamera()
	switch evt.Keysym.Scancode {
	case sdl.SCANCODE_W:
		cam.Translate(cam.GetLookForward())
	case sdl.SCANCODE_A:
		cam.Translate(cam.GetLookLeft())
	case sdl.SCANCODE_S:
		cam.Translate(cam.GetLookBack())
	case sdl.SCANCODE_D:
		cam.Translate(cam.GetLookRight())
	default:
		return
	}
}

func (app *Application) PostEventActions() {
	w, h := app.win.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	err := app.plane.Render()
	if err != nil {
		log.Warnf("plane render error: %v", err)
	}

	app.win.GLSwap()

	for glErr := gl.GetError(); glErr != gl.NO_ERROR; glErr = gl.GetError() {
		log.Warnf("OpenGL error: %v", glErr)
	}
}

func (app *Application) Quit() {
	if err := app.win.Destroy(); err != nil {
		log.Fatal(err)
	}
	sdl.Quit()
}
