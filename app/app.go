package app

import (
	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/game"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/oldworld"
	"github.com/kroppt/voxels/physics"
	"github.com/kroppt/voxels/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win     *sdl.Window
	world   *oldworld.World
	running bool
	m1held  bool
	game    *game.Game
}

func New(win *sdl.Window) (*Application, error) {
	return &Application{
		win:   win,
		world: oldworld.New(),
		game:  game.New(game.OsTimeNow),
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
	// key := block.Pos.GetChunkPos(world.ChunkSize)
	// log.Debug(app.world.ChunksLoaded[key].Chunk)
	if evt.Y < 0 {
		sw := util.Start()
		app.world.RemoveVoxel(block.Pos)
		sw.StopRecordAverage("remove voxel")
	} else {
		sw := util.Start()
		app.world.SetVoxel(&oldworld.Voxel{
			Pos:   block.Pos,
			Btype: oldworld.Light,
		})
		sw.StopRecordAverage("set voxel")
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

	vel := physics.Vel{}
	if keys[sdl.SCANCODE_W] == sdl.PRESSED {
		look := cam.GetLookForward()
		vel.Vec3 = vel.Add(&look)
	}
	if keys[sdl.SCANCODE_S] == sdl.PRESSED {
		look := cam.GetLookBack()
		vel.Vec3 = vel.Add(&look)
	}
	if keys[sdl.SCANCODE_A] == sdl.PRESSED {
		look := cam.GetLookLeft()
		vel.Vec3 = vel.Add(&look)
	}
	if keys[sdl.SCANCODE_D] == sdl.PRESSED {
		look := cam.GetLookRight()
		vel.Vec3 = vel.Add(&look)
	}
	if keys[sdl.SCANCODE_SPACE] == sdl.PRESSED {
		look := glm.Vec3{0.0, 1.0, 0.0}
		vel.Vec3 = vel.Add(&look)
	}
	if keys[sdl.SCANCODE_LSHIFT] == sdl.PRESSED {
		look := glm.Vec3{0.0, -1.0, 0.0}
		vel.Vec3 = vel.Add(&look)
	}

	if vel.Len() == 0 {
		return nil
	}

	speed := float32(10.0)
	vel.Vec3 = vel.Vec3.Normalized()
	vel.Vec3 = vel.Vec3.Mul(speed)
	dt := app.game.GetTickDuration()
	pos := vel.AsPosition(dt)
	cam.Translate(&pos.Vec3)

	return nil
}

func (app *Application) PostEventActions() {
	app.game.NextTick()
	app.pollKeyboard()
	w, h := app.win.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	sw := util.Start()
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
