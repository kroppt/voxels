package graphics

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/kroppt/voxels/util"
	"github.com/veandco/go-sdl2/sdl"
)

type core struct {
	ubo            *gfx.BufferObject
	textureMap     *gfx.CubeMap
	crosshair      *CrosshairObject
	selectedVoxel  SelectedVoxel
	selected       bool
	selectionFrame *glObject
	window         *sdl.Window
	settingsRepo   settings.Interface
	loadedChunks   map[chunk.ChunkCoordinate]*glObject
	viewableChunks map[chunk.ChunkCoordinate]struct{}
}

// ErrRenderDriver indicates that SDL failed to enable the OpenGL render driver.
const ErrRenderDriver log.ConstErr = "failed to set opengl render driver hint"

func (c *core) createWindow(title string) error {
	runtime.LockOSThread()

	width, height := c.settingsRepo.GetResolution()
	if width == 0 || height == 0 {
		width = 1280
		height = 720
	}

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

	ubo := gfx.NewBufferObject()
	var mat mgl.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)

	c.textureMap = loadSpriteSheet("sprite_sheet.png")

	c.crosshair, err = NewCrosshairObject(float32(c.settingsRepo.GetCrosshairLength()), float32(width)/float32(height))
	if err != nil {
		return fmt.Errorf("failed to make crosshair: %v", err)
	}

	c.selectionFrame, err = newFrameObject()
	if err != nil {
		return fmt.Errorf("failed to make selection frame: %v", err)
	}

	c.window = window
	c.ubo = ubo
	return nil
}

func loadSpriteSheet(fileName string) *gfx.CubeMap {
	// TODO get data without texture
	sprites, err := gfx.NewTextureFromFile(fileName)
	if err != nil {
		panic("failed to load sprite sheet")
	}
	sprytes := sprites.GetData()
	// TODO make fancy file format with meta data
	w := int32(16)
	h := sprites.GetHeight()
	layers := h / w
	texAtlas, err := gfx.NewCubeMap(w, layers, sprytes, gl.RGBA, 4, 4)
	if err != nil {
		panic("failed to create 3d texture")
	}
	texAtlas.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	texAtlas.SetParameter(gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	return &texAtlas
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
	c.ubo.Destroy()
	c.textureMap.Destroy()
	c.crosshair.Destroy()
	c.selectionFrame.destroy()
	sdl.Quit()
	return err
}

func (c *core) getUpdatedProjMatrix() mgl.Mat4 {
	fovRad := mgl.DegToRad(c.settingsRepo.GetFOV())
	near := c.settingsRepo.GetNear()
	far := c.settingsRepo.GetFar()
	width, height := c.settingsRepo.GetResolution()
	aspect := float64(width) / float64(height)
	return mgl.Perspective(fovRad, aspect, near, far)
}

func (c *core) updateView(viewableChunks map[chunk.ChunkCoordinate]struct{}, view mgl.Mat4, selectedVoxel SelectedVoxel, selected bool) {
	c.viewableChunks = viewableChunks

	if selected && c.selected && selectedVoxel != c.selectedVoxel {
		c.selectionFrame.setData([]float32{selectedVoxel.X, selectedVoxel.Y, selectedVoxel.Z, selectedVoxel.Vbits, 0})
		c.selectedVoxel = selectedVoxel
	}
	c.selected = selected

	proj := c.getUpdatedProjMatrix()
	err := c.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		panic(fmt.Sprintf("failed to upload camera view to ubo: %v", err))
	}
	err = c.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(view)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		panic(fmt.Sprintf("failed to upload camera proj to ubo: %v", err))
	}
}

func (c *core) loadChunk(chunk chunk.Chunk) {
	if _, ok := c.loadedChunks[chunk.Position()]; ok {
		panic("attempting to load over an already-loaded chunk")
	}

	chunkObj, err := newChunkObject()
	if err != nil {
		panic(err)
	}
	chunkObj.setData(chunk.GetFlatData())

	c.loadedChunks[chunk.Position()] = chunkObj
}

func (c *core) unloadChunk(key chunk.ChunkCoordinate) {
	if _, ok := c.loadedChunks[key]; !ok {
		panic("attempting to unload a chunk that is not loaded")
	}
	c.loadedChunks[key].destroy()
	delete(c.loadedChunks, key)
}

func (c *core) render() {
	w, h := c.window.GetSize()
	gl.Viewport(0, 0, w, h)

	// clear with black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	sw := util.Start()

	c.textureMap.Bind()
	for _, chunkObj := range c.loadedChunks {
		chunkObj.render()
	}
	c.textureMap.Unbind()

	gl.Disable(gl.DEPTH_TEST)
	if c.selected {
		gl.LineWidth(2.0)
		c.selectionFrame.render()
	}
	gl.LineWidth(float32(c.settingsRepo.GetCrosshairThickness()))
	c.crosshair.Render()
	gl.Enable(gl.DEPTH_TEST)

	sw.StopRecordAverage("Total World Render")
	c.window.GLSwap()

	for glErr := gl.GetError(); glErr != gl.NO_ERROR; glErr = gl.GetError() {
		log.Warnf("OpenGL error: %v", glErr)
	}
}
