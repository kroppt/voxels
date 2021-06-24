package world

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/util"
	"github.com/kroppt/voxels/voxgl"
)

// World tracks the camera and its renderable chunks.
type World struct {
	ubo               *gfx.BufferObject
	cam               *Camera
	chunksLoaded      map[ChunkPos]*LoadedChunk
	chunkExpectRender map[ChunkPos]struct{}
	chunkExpect       map[ChunkPos]struct{}
	chunkSaving       map[ChunkPos]struct{}
	chunkLoading      map[ChunkPos]struct{}
	currChunk         ChunkPos
	chunkChan         chan *Chunk
	saved             chan ChunkPos
	loaded            chan ChunkPos
	cubeMap           *gfx.CubeMap
	gen               Generator
	cache             *Cache
	cacheLock         sync.Mutex
	cancel            bool
	selectedVoxel     *voxgl.Object
	selected          bool
	crosshair         *voxgl.Crosshair
}

type LoadedChunk struct {
	chunk    *Chunk
	modified bool
	doRender bool
	relight  bool
}

type AdjMod struct {
	off     VoxelPos
	adjDiff AdjacentMask
	add     bool
}

const ChunkSize = 8
const chunkRenderDist = 3

// chunkRenderBuffer gaurantees a minimum radius of area of effect operations
const chunkRenderBuffer = 1
const cacheThreshold = 10
const selectionDist = 10

// New returns a new world.World.
func New() *World {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	cam := NewCameraDefault()
	crosshair, err := voxgl.NewCrosshair(0.03, cam.aspect)
	if err != nil {
		panic(fmt.Sprintf("failed to make crosshair: %v", err))
	}

	world := &World{
		ubo:       ubo,
		cam:       cam,
		chunkChan: make(chan *Chunk),
		saved:     make(chan ChunkPos),
		loaded:    make(chan ChunkPos),
		gen:       FlatWorldGenerator{},
		crosshair: crosshair,
		cacheLock: sync.Mutex{},
	}

	cam.SetPosition(&glm.Vec3{0.5, 7.5, 2})
	// cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})

	cache, err := NewCache("world_meta", "world_data", cacheThreshold)
	if err != nil {
		panic(fmt.Sprint(err))
	}
	world.cache = cache

	rand.Seed(time.Now().UnixNano())

	world.chunksLoaded = map[ChunkPos]*LoadedChunk{}
	world.chunkExpectRender = make(map[ChunkPos]struct{})
	world.chunkExpect = make(map[ChunkPos]struct{})
	world.chunkSaving = make(map[ChunkPos]struct{})
	world.chunkLoading = make(map[ChunkPos]struct{})

	world.selectedVoxel, err = voxgl.NewFrame(nil)
	if err != nil {
		panic(fmt.Sprintf("error creating NewFrame: %v", err))
	}

	world.cubeMap = loadSpriteSheet("sprite_sheet.png")

	world.update()
	return world
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

// FindLookAtVoxel determines which voxel is being looked at. It returns the
// block, distance to the block, and whether the block was found.
func (w *World) FindLookAtVoxel() (block *Voxel, dist float32, found bool) {
	var bestVox *Voxel
	var bestDist float32
	for _, loadedCh := range w.chunksLoaded {
		if !loadedCh.doRender {
			continue
		}
		vox, dist, hit := loadedCh.chunk.root.FindClosestIntersect(w.cam)
		if hit && (dist < bestDist || bestVox == nil) {
			bestVox = vox
			bestDist = dist
		}
	}
	return bestVox, bestDist, bestVox != nil
}

func (w *World) updateSelectedVoxel() {
	sw := util.Start()
	v, dist, found := w.FindLookAtVoxel()
	sw.StopRecordAverage("FindLookAtVoxel")
	if !found || dist > selectionDist {
		w.selected = false
		return
	} else {
		w.selected = true
	}
	key := v.Pos.GetChunkPos(ChunkSize)
	loadedCh, ok := w.chunksLoaded[key]
	if !ok {
		panic("expected look at voxel to be in a loaded chunk")
	}
	vox := loadedCh.chunk.GetVoxelFromFlatData(v.Pos)
	pos := vox.Pos.AsVec3()
	w.selectedVoxel.SetData([]float32{pos.X(), pos.Y(), pos.Z(), float32(vox.GetVbits())})
}

// SetVoxel updates a voxel's variables in the world if the chunk
// that it would belong to is currently loaded.
func (w *World) SetVoxel(v *Voxel) {
	// TODO setting a voxel to air is *similar* to removing, but not the same
	// -> what would setting a voxel to air do here? should this be prevented?
	key := v.Pos.GetChunkPos(ChunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	loadedCh, ok := w.chunksLoaded[key]
	if !ok {
		panic(fmt.Sprintf("tried to set a voxel (%v) that belongs to an unrendered chunk", v.Pos))
	}
	chunk := loadedCh.chunk
	target := chunk.GetVoxelFromFlatData(v.Pos)
	v.AdjMask = target.AdjMask
	chunk.SetVoxel(v)
	mods := []AdjMod{
		{VoxelPos{-1, 0, 0}, AdjacentRight, true},
		{VoxelPos{1, 0, 0}, AdjacentLeft, true},
		{VoxelPos{0, -1, 0}, AdjacentTop, true},
		{VoxelPos{0, 1, 0}, AdjacentBottom, true},
		{VoxelPos{0, 0, -1}, AdjacentBack, true},
		{VoxelPos{0, 0, 1}, AdjacentFront, true},
	}
	for _, mod := range mods {
		mod.off = mod.off.Add(v.Pos)
		k := mod.off.GetChunkPos(ChunkSize)
		ch, ok := w.chunksLoaded[k]
		if !ok {
			panic("not currently handling pending mods to unloaded chunks")
		} else {
			ch.chunk.AddAdjacency(mod.off, mod.adjDiff)
			ch.modified = true
		}
	}

	if v.Btype == Light {
		w.lightFrom(chunk, v.Pos, MaxLightValue, v.Pos)
	} else if target.Btype != Light {
		// TODO setting a light block to something else, update lighting
		w.lightRemoveFrom(chunk, v.Pos, v.Pos)

	}
}

func (w *World) relightFromNeighbors(v VoxelPos) {
	key := v.GetChunkPos(ChunkSize)
	chunk := w.chunksLoaded[key].chunk
	srcMap := chunk.lightRefs[v]
	uniques := make(map[VoxelPos]uint32)
	for src, val := range srcMap {
		uniques[src] = val
	}
	for _, mod := range lightMods {
		offP := v.Add(mod.off)
		offKey := offP.GetChunkPos(ChunkSize)
		if loadedCh, ok := w.chunksLoaded[offKey]; ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			if adjBlock.Btype == Air {
				offSrcMap := offCh.lightRefs[offP]
				for offSrc, offVal := range offSrcMap {
					uniques[offSrc] = offVal
				}
			}
		}
	}
	for uniqueSrc := range uniques {
		srcChKey := uniqueSrc.GetChunkPos(ChunkSize)
		srcCh := w.chunksLoaded[srcChKey].chunk
		// TODO have custom light value per light so its not hacked in here
		w.lightFrom(srcCh, uniqueSrc, MaxLightValue, uniqueSrc)
	}
}

func (w *World) RemoveVoxel(v VoxelPos) {
	key := v.GetChunkPos(ChunkSize)
	loadedCh, ok := w.chunksLoaded[key]
	if !ok {
		panic(fmt.Sprintf("tried to remove a voxel (%v) that belongs to an unloaded chunk", v))
	}
	chunk := loadedCh.chunk
	oldType := chunk.GetVoxelFromFlatData(v).Btype
	// log.Debugf("removing voxel in chunk %v", key)
	// TODO setting voxel to air and removing node from octree are tied
	// together in a "delete voxel" operation, this should be more well defined
	// inside chunk.go in a single function call
	chunk.SetVoxel(&Voxel{
		Pos:     v,
		Btype:   Air,
		AdjMask: AdjacentNone,
	})
	chunk.root, _ = chunk.root.Remove(v)
	mods := []AdjMod{
		{VoxelPos{-1, 0, 0}, AdjacentRight, false},
		{VoxelPos{1, 0, 0}, AdjacentLeft, false},
		{VoxelPos{0, -1, 0}, AdjacentTop, false},
		{VoxelPos{0, 1, 0}, AdjacentBottom, false},
		{VoxelPos{0, 0, -1}, AdjacentBack, false},
		{VoxelPos{0, 0, 1}, AdjacentFront, false},
	}
	for _, mod := range mods {
		mod.off = mod.off.Add(v)
		k := mod.off.GetChunkPos(ChunkSize)
		ch, ok := w.chunksLoaded[k]
		if !ok {
			panic("not currently handling pending mods to unloaded chunks")
		} else {
			ch.chunk.RemoveAdjacency(mod.off, mod.adjDiff)
			ch.modified = true
		}
	}
	if oldType == Light {
		w.lightRemoveFrom(chunk, v, v)
	}
	// TODO can this be too slow?
	// this is required because lesser-light sources might remain
	// if the stronger light source was removed, so the book keeping
	// needs to be calculated
	w.relightFromNeighbors(v)
}

// GetCamera returns a reference to the camera.
func (w *World) GetCamera() *Camera {
	return w.cam
}

// TODO Create UBO extraction
func (w *World) updateUBO() error {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		return err
	}
	proj := cam.GetProjMat()
	err = w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(view)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		return err
	}
	return nil
}

// TODO when/where will this be used? chunk constructor?
func (w *World) updateAllLights(ch *Chunk) {
	for lightPos := range ch.lights {
		w.lightFrom(ch, lightPos, MaxLightValue, lightPos)
	}
}

var lightMods = [6]struct {
	off  VoxelPos
	face LightMask
}{
	{VoxelPos{-1, 0, 0}, LightRight},
	{VoxelPos{1, 0, 0}, LightLeft},
	{VoxelPos{0, -1, 0}, LightTop},
	{VoxelPos{0, 1, 0}, LightBottom},
	{VoxelPos{0, 0, -1}, LightBack},
	{VoxelPos{0, 0, 1}, LightFront},
}

func (w *World) lightRemoveFrom(c *Chunk, p VoxelPos, src VoxelPos) {
	if !c.IsWithinChunk(p) {
		panic("improper usage: p outside chunk")
	}
	if c.HasSource(p, src) {
		c.DeleteSource(p, src)
	} else {
		return
	}

	for _, mod := range lightMods {
		offP := p.Add(mod.off)
		offKey := offP.GetChunkPos(ChunkSize)
		if loadedCh, ok := w.chunksLoaded[offKey]; ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			if adjBlock.Btype == Air {
				w.lightRemoveFrom(offCh, offP, src)
			} else {
				_, secondLargest, found := c.GetBrightestSource(p)
				if found {
					adjBlock.SetLightValue(secondLargest, mod.face)
				} else {
					adjBlock.SetLightValue(0, mod.face)
				}
				offCh.SetVoxel(&adjBlock)
			}
		}
	}

}

type LightMod struct {
	pos   VoxelPos
	src   VoxelPos
	value uint32
	face  LightMask
}

func (w *World) lightFrom(c *Chunk, p VoxelPos, value uint32, src VoxelPos) {
	if value < 0 || value > MaxLightValue || !c.IsWithinChunk(p) {
		panic("improper usage: bad lighting value or p outside chunk")
	}
	if value == 0 || (c.HasSource(p, src) && c.GetSourceValue(p, src) > value) {
		return
	} else {
		c.SetSource(p, src, value)
	}
	if p == src {
		v := c.GetVoxelFromFlatData(p)
		v.SetLightValue(MaxLightValue, LightLeft)
		v.SetLightValue(MaxLightValue, LightRight)
		v.SetLightValue(MaxLightValue, LightBack)
		v.SetLightValue(MaxLightValue, LightFront)
		v.SetLightValue(MaxLightValue, LightTop)
		v.SetLightValue(MaxLightValue, LightBottom)
		c.SetVoxel(&v)
	}
	// TODO refactor
	mods := []LightMod{
		{VoxelPos{-1, 0, 0}, src, value, LightRight},
		{VoxelPos{1, 0, 0}, src, value, LightLeft},
		{VoxelPos{0, -1, 0}, src, value, LightTop},
		{VoxelPos{0, 1, 0}, src, value, LightBottom},
		{VoxelPos{0, 0, -1}, src, value, LightBack},
		{VoxelPos{0, 0, 1}, src, value, LightFront},
	}
	for _, mod := range mods {
		mod.pos = p.Add(mod.pos)
		offKey := mod.pos.GetChunkPos(ChunkSize)
		if loadedCh, ok := w.chunksLoaded[offKey]; ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(mod.pos)
			if adjBlock.Btype == Air {
				// continue search down open path
				w.lightFrom(offCh, mod.pos, value-1, src)
			} else {
				// path is blocked off, apply lighting value to the face *if brighter*
				if adjBlock.GetLightValue(mod.face) < value {
					adjBlock.SetLightValue(value, mod.face)
					offCh.SetVoxel(&adjBlock)
				}
			}
		}
	}
}

func (w *World) isWithinRenderDist(key ChunkPos) bool {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)
	return rng.Contains(key)
}

func (w *World) hasSurroundingChunks(key ChunkPos) bool {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	// loadedRng handles the edge case for buffered, invisible
	// chunks that aren't expected to have chunks beyond them
	// loaded before they run their lighting algorithm.
	loadedRng := currChunk.GetSurroundings(chunkRenderDist + chunkRenderBuffer)
	localBufRng := key.GetSurroundings(chunkRenderBuffer)
	allHere := true
	localBufRng.ForEach(func(pos ChunkPos) {
		if _, ok := w.chunksLoaded[pos]; !ok && loadedRng.Contains(pos) {
			allHere = false
		}
	})
	return allHere
}

// receiveExpectedAsync reads loaded chunks off the chunk channel, and only adds them
// to the collection of chunks to be rendered if they are still expected (not late)
func (w *World) receiveExpectedAsync() {
	for {
		select {
		case ch := <-w.chunkChan:
			if _, ok := w.chunkExpect[ch.Pos]; ok {
				// the chunk has arrived and we expected it
				// give the chunk its object
				objs, err := voxgl.NewChunkObject(nil)
				if err != nil {
					panic(fmt.Sprint(err))
				}
				ch.SetObjs(objs)

				w.chunksLoaded[ch.Pos] = &LoadedChunk{
					chunk:    ch,
					modified: false,
					relight:  true,
					doRender: w.isWithinRenderDist(ch.Pos),
				}

				// check waiting list to see if this arrival
				// completed a different chunk's surroundings
				// TODO only check within buffer-range of arriving chunk
				for key, loadedCh := range w.chunksLoaded {
					if loadedCh.relight && w.hasSurroundingChunks(key) {
						w.updateAllLights(loadedCh.chunk)
						loadedCh.relight = false
					}
				}
			}
		default:
			return
		}
	}
}

// requestChunk puts in an asynchronous request to the cache for a chunk
// at a particular ChunkPos to be loaded, which will complete at *some
// point in the future*, quite possibly when it's no longer expected.
// When the async call completes, the loaded chunk will be placed on the
// chunk channel, and the chunk's key will be placed on the loaded channel.
// A request to any chunk that is already loaded, loading, or saving will
// be ignored.
func (w *World) requestChunk(key ChunkPos) {
	_, loaded := w.chunksLoaded[key]
	_, loading := w.chunkLoading[key]
	_, saving := w.chunkSaving[key]
	if loaded || loading || saving {
		return
	}
	w.chunkLoading[key] = struct{}{}
	go func(key ChunkPos) {
		// check cache for a saved chunk
		w.cacheLock.Lock()
		if w.cancel {
			w.cacheLock.Unlock()
			return
		}
		chunk, loaded := w.cache.Load(key)
		w.cacheLock.Unlock()
		if !loaded {
			chunk = NewChunk(ChunkSize, key, w.gen)
		}

		// TODO wait for surrounding extra chunks to load for lighting
		// only if this is a renderable chunk

		// TODO switching these lines changes things
		// should these two be tied?
		w.chunkChan <- chunk
		w.loaded <- key
	}(key)
}

// requestExpectedChunks attempts to request every expected chunk
// Note that in requestChunk, it will ignore requests to loaded,
// loading or saving chunks.
func (w *World) requestExpectedChunks() {
	for key := range w.chunkExpect {
		w.requestChunk(key)
	}
}

// checkSavingStatus reads chunk keys off the saved channel to indicate
// that a particular chunk has finished saving. More importantly, if the
// chunk is still within render distance, it marks the chunk as expected and
// requests the chunk to be loaded.
func (w *World) checkSavingStatus() {
	for {
		select {
		case key := <-w.saved:
			delete(w.chunkSaving, key)
		default:
			return
		}
	}
}

// checkLoadingStatus reads chunk keys off the loaded channel
// to indicate via the chunkLoading map that a particular chunk is no longer
// being loaded.
func (w *World) checkLoadingStatus() {
	for {
		select {
		case key := <-w.loaded:
			delete(w.chunkLoading, key)
		default:
			return
		}
	}
}

// Expected chunks are those that are within the render distance from the
// chunk that the camera is currently in, and are also not in the process
// of being saved.
func (w *World) updateExpectedChunks() {
	// slightly larger range for chunks that should be loaded
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist + chunkRenderBuffer)
	w.chunkExpect = make(map[ChunkPos]struct{})
	rng.ForEach(func(pos ChunkPos) {
		if _, saving := w.chunkSaving[pos]; !saving {
			// only expect chunks that are not in the process of saving
			w.chunkExpect[pos] = struct{}{}
		}
	})
}

// evictUnexpectedChunks checks all loaded chunks that are being rendered
// and removes any that are no longer expected. If they were modified, it
// puts in an asynchronous request for the chunk to be saved, marking its key
// as in the process of saving, and indicating completion on the saved channel
// when it is done.
func (w *World) evictUnexpectedChunks() {
	for key, loadedCh := range w.chunksLoaded {
		if _, ok := w.chunkExpect[key]; !ok {
			w.chunksLoaded[key].chunk.Destroy()
			delete(w.chunksLoaded, key)
			if loadedCh.modified {
				w.chunkSaving[key] = struct{}{}
				go func(key ChunkPos, ch *Chunk) {
					w.cacheLock.Lock()
					if w.cancel {
						w.cacheLock.Unlock()
						return
					}
					w.cache.Save(ch)
					w.cacheLock.Unlock()
					w.saved <- key
				}(key, loadedCh.chunk)
			}
		} else {
			// TODO rename evictUnexpectedChunks to account for this logic?
			shouldRender := w.isWithinRenderDist(key)
			if !loadedCh.doRender && shouldRender {
				// this chunk was a buffer chunk and has moved into render dist
				loadedCh.relight = true
			}
			loadedCh.doRender = shouldRender
		}
	}
}

func (w *World) update() {
	w.updateExpectedChunks()
	w.evictUnexpectedChunks()
	w.requestExpectedChunks()
}

// Render renders the chunks of the world in OpenGL.
func (w *World) Render() error {
	sw := util.Start()
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	if currChunk != w.currChunk {
		w.currChunk = currChunk
		w.update()
	}

	w.checkSavingStatus()
	w.checkLoadingStatus()
	w.receiveExpectedAsync()
	sw.StopRecordAverage("total update logic")

	if w.cam.IsDirty() {
		err := w.updateUBO()
		if err != nil {
			return err
		}
		w.cam.Clean()
	}

	for _, loadedCh := range w.chunksLoaded {
		if loadedCh.doRender {
			w.cubeMap.Bind()
			loadedCh.chunk.Render(w.cam)
			w.cubeMap.Unbind()
		}
	}
	gl.LineWidth(2)
	gl.Disable(gl.DEPTH_TEST)
	w.updateSelectedVoxel()
	if w.selected {
		w.selectedVoxel.Render()
	}
	w.crosshair.Render()
	gl.Enable(gl.DEPTH_TEST)

	return nil
}

func (w *World) saveRoutine() {
	w.cancel = true
	w.cacheLock.Lock()
	// TODO chunksloaded, not just to render, because invisible loaded chunks could be modified too
	for _, ch := range w.chunksLoaded {
		if ch.modified {
			w.cache.Save(ch.chunk)
		}
	}
	w.cache.Destroy()
	w.cacheLock.Unlock()
}

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
	w.saveRoutine()
}
