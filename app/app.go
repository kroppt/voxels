package app

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/world"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win     *sdl.Window
	plane   *world.Plane
	running bool
	m1held  bool
}

func New(win *sdl.Window) (*Application, error) {
	// 11 x 11 x 11
	x := world.Range{Min: -5, Max: 5}
	y := world.Range{Min: -5, Max: 5}
	z := world.Range{Min: -5, Max: 5}
	plane, err := world.NewPlane(NewPlaneRenderer(), x, y, z)
	if err != nil {
		return nil, fmt.Errorf("could not create plane: %v", err)
	}

	return &Application{
		win:   win,
		plane: plane,
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
func (app *Application) HandleSdlEvent(e sdl.Event) error {
	switch evt := e.(type) {
	case *sdl.QuitEvent:
		app.handleQuitEvent(evt)
	case *sdl.MouseButtonEvent:
		app.handleMouseButtonEvent(evt)
	case *sdl.MouseMotionEvent:
		err := app.handleMouseMotionEvent(evt)
		if err != nil {
			return err
		}
	case *sdl.MouseWheelEvent:
		//app.handleMouseWheelEvent(evt)
	case *sdl.WindowEvent:
		//app.handleWindowEvent(evt)
	case *sdl.SysWMEvent:
		//app.handleSysWMEvent(evt)
	case *sdl.KeyboardEvent:
		err := app.handleKeyboardEvent(evt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) handleQuitEvent(evt *sdl.QuitEvent) {
	app.running = false
}

func (app *Application) handleMouseButtonEvent(evt *sdl.MouseButtonEvent) {
	if evt.State == sdl.PRESSED && !app.m1held {
		app.m1held = true
	} else if evt.State == sdl.RELEASED {
		app.m1held = false
	}
}

func (app *Application) handleMouseMotionEvent(evt *sdl.MouseMotionEvent) error {
	if !app.m1held {
		return nil
	}
	cam := app.plane.GetCamera()
	speed := float32(0.1)
	// use x component to rotate around Y axis
	cam.Rotate(mgl.Vec3{0.0, 1.0, 0.0}, speed*float32(evt.XRel))
	// use y component to rotate around the axis that goes through your ears
	cam.Rotate(cam.GetLookRight(), speed*float32(evt.YRel))
	err := app.plane.GetRenderer().UpdateView()
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) handleKeyboardEvent(evt *sdl.KeyboardEvent) error {
	// TODO check/handle multiple keys pressed simultaneously for keyboards
	// that don't rapidly spam both keys
	if evt.State != sdl.PRESSED {
		return nil
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
	case sdl.SCANCODE_SPACE:
		cam.Translate(mgl.Vec3{0.0, 1.0, 0.0})
	case sdl.SCANCODE_LSHIFT:
		cam.Translate(mgl.Vec3{0.0, -1.0, 0.0})
	default:
		return nil
	}
	err := app.plane.GetRenderer().UpdateView()
	if err != nil {
		return err
	}
	return nil
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
	app.plane.Destroy()
	if err := app.win.Destroy(); err != nil {
		log.Fatal(err)
	}
	sdl.Quit()
}
