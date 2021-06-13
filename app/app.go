package app

import (
	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/util"
	"github.com/kroppt/voxels/world"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win     *sdl.Window
	world   *world.World
	running bool
	m1held  bool
}

func New(win *sdl.Window) (*Application, error) {
	return &Application{
		win:   win,
		world: world.New(),
	}, nil
}

func (app *Application) Start() {
	util.SetMetricsEnabled(true)
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
		app.handleMouseWheelEvent(evt)
	case *sdl.WindowEvent:
		//app.handleWindowEvent(evt)
	case *sdl.SysWMEvent:
		//app.handleSysWMEvent(evt)
	case *sdl.KeyboardEvent:
		app.handleKeyboardEvent(evt)
	}
	return nil
}

func (app *Application) handleMouseWheelEvent(evt *sdl.MouseWheelEvent) {
	block, _, ok := app.world.FindLookAtVoxel()
	if !ok {
		return
	}
	if evt.Y < 0 {
		app.world.RemoveVoxel(block.Pos)
	} else {
		app.world.SetVoxel(&world.Voxel{
			Pos:   block.Pos,
			Btype: world.Labeled,
		})
	}
}

func (app *Application) handleQuitEvent(evt *sdl.QuitEvent) {
	util.LogMetrics()
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
	cam := app.world.GetCamera()
	speed := float32(0.1)
	// use x component to rotate around Y axis
	cam.Rotate(&glm.Vec3{0.0, 1.0, 0.0}, speed*float32(evt.XRel))
	// use y component to rotate around the axis that goes through your ears
	lookRight := cam.GetLookRight()
	cam.Rotate(&lookRight, speed*float32(evt.YRel))
	return nil
}

func (app *Application) handleKeyboardEvent(evt *sdl.KeyboardEvent) {
}

func (app *Application) pollKeyboard() error {
	cam := app.world.GetCamera()
	keys := sdl.GetKeyboardState()
	speed := float32(0.07)
	if keys[sdl.SCANCODE_W] == sdl.PRESSED {
		look := cam.GetLookForward()
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	if keys[sdl.SCANCODE_S] == sdl.PRESSED {
		look := cam.GetLookBack()
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	if keys[sdl.SCANCODE_A] == sdl.PRESSED {
		look := cam.GetLookLeft()
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	if keys[sdl.SCANCODE_D] == sdl.PRESSED {
		look := cam.GetLookRight()
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	if keys[sdl.SCANCODE_SPACE] == sdl.PRESSED {
		look := glm.Vec3{0.0, 1.0, 0.0}
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	if keys[sdl.SCANCODE_LSHIFT] == sdl.PRESSED {
		look := glm.Vec3{0.0, -1.0, 0.0}
		lookSpeed := look.Mul(speed)
		cam.Translate(&lookSpeed)
	}
	return nil
}

var Block *world.Voxel

func (app *Application) PostEventActions() {
	app.pollKeyboard()
	sw := util.Start()
	Block, _, _ = app.world.FindLookAtVoxel()
	// if found {
	// 	cam := app.world.GetCamera()
	// 	eye := cam.GetPosition()
	// 	dir := cam.GetLookForward()
	// 	log.Debugf("Block: %v, dist: %v, pos: %v, look: %v", Block.Pos, dist, eye, dir)
	// }
	sw.StopRecordAverage("Intersect")
	w, h := app.win.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	sw = util.Start()
	app.world.Render()
	sw.StopRecordAverage("Total World Render")
	app.win.GLSwap()

	for glErr := gl.GetError(); glErr != gl.NO_ERROR; glErr = gl.GetError() {
		log.Warnf("OpenGL error: %v", glErr)
	}
}

func (app *Application) Quit() {
	app.world.Destroy()
	if err := app.win.Destroy(); err != nil {
		log.Fatal(err)
	}
	sdl.Quit()
}
