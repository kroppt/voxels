package app

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/log"
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	win     *sdl.Window
	cube    *sampleCube
	running bool
}

func New(win *sdl.Window) (*Application, error) {
	cube, err := newSampleCube()
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
	}
}

func (app *Application) Render() {
	app.cube.Render()
}

func (app *Application) handleQuitEvent(evt *sdl.QuitEvent) {
	app.running = false
}

func (app *Application) PostEventActions() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	// clear with black
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	app.Render()

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

type sampleCube struct {
	prog gfx.Program
	buf  *gfx.VAO
}

// close / far (red)
// top / bottom (blue)
// left / right (green)
var fbl = [6]float32{-0.5, -0.5, 0.7, 0.8, 0.8, 0.2}
var fbr = [6]float32{0.5, -0.5, 0.7, 0.8, 0.8, 0.8}
var ftl = [6]float32{-0.5, 0.5, 0.7, 0.8, 0.2, 0.2}
var ftr = [6]float32{0.5, 0.5, 0.7, 0.8, 0.2, 0.8}

var cbl = [6]float32{-0.7, -0.7, 0.3, 0.2, 0.8, 0.2}
var cbr = [6]float32{0.3, -0.7, 0.3, 0.2, 0.8, 0.8}
var ctl = [6]float32{-0.7, 0.3, 0.3, 0.2, 0.2, 0.2}
var ctr = [6]float32{0.3, 0.3, 0.3, 0.2, 0.2, 0.8}

func newSampleCube() (*sampleCube, error) {
	var err error

	v1, err := gfx.NewShader(gfx.SampleCubeVertex, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	f1, err := gfx.NewShader(gfx.SampleCubeFragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(v1, f1)
	if err != nil {
		return nil, err
	}

	buf := gfx.NewVAO(gl.TRIANGLES, []int32{3, 3})

	vertices := [][6]float32{
		// far face
		fbl, ftl, fbr,
		ftl, fbr, ftr,

		// left face
		ftl, ctl, fbl,
		ctl, fbl, cbl,

		// top face
		ftl, ctl, ftr,
		ctl, ftr, ctr,

		// right face
		fbr, ftr, cbr,
		ftr, cbr, ctr,

		// bottom face
		cbl, cbr, fbl,
		cbr, fbl, fbr,

		// close face
		ctl, ctr, cbr,
		cbl, ctl, cbr,
	}

	points := []float32{}
	for _, v := range vertices {
		points = append(points, v[:]...)
	}

	err = buf.Load(points, gl.STATIC_DRAW)
	if err != nil {
		return nil, err
	}

	return &sampleCube{
		prog: prog,
		buf:  buf,
	}, nil
}

func (sc *sampleCube) Render() {
	sc.prog.Bind()
	sc.buf.Draw()
	sc.prog.Unbind()
}

func (sc *sampleCube) Destroy() {
	sc.prog.Destroy()
	sc.buf.Destroy()
}
