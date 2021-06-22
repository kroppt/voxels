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
	ubo           *gfx.BufferObject
	cam           *Camera
	chunks        map[ChunkPos]*Chunk
	chunkExpect   map[ChunkPos]struct{}
	chunkSaving   map[ChunkPos]struct{}
	chunkLoading  map[ChunkPos]struct{}
	pendingFuncs  map[ChunkPos][]func(*Chunk)
	currChunk     ChunkPos
	chunkChan     chan *Chunk
	saved         chan ChunkPos
	loaded        chan ChunkPos
	cubeMap       *gfx.CubeMap
	gen           Generator
	cache         *Cache
	cacheLock     sync.Mutex
	cancel        bool
	selectedVoxel *voxgl.Object
	selected      bool
	crosshair     *voxgl.Crosshair
}

type AdjMod struct {
	off     VoxelPos
	adjDiff AdjacentMask
	add     bool
}

const ChunkSize = 4
const chunkRenderDist = 3
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
	world.chunks = make(map[ChunkPos]*Chunk)
	world.chunkExpect = make(map[ChunkPos]struct{})
	world.chunkLoading = make(map[ChunkPos]struct{})
	world.chunkSaving = make(map[ChunkPos]struct{})
	world.pendingFuncs = make(map[ChunkPos][]func(*Chunk))
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
	for _, chunk := range w.chunks {
		vox, dist, hit := chunk.root.FindClosestIntersect(w.cam)
		if hit && (dist < bestDist || bestVox == nil) {
			bestVox = vox
			bestDist = dist
		}
	}
	return bestVox, bestDist, bestVox != nil
}

func (w *World) updateSelectedVoxel() {
	v, dist, found := w.FindLookAtVoxel()
	if !found || dist > selectionDist {
		w.selected = false
		return
	} else {
		w.selected = true
	}
	key := v.Pos.GetChunkPos(ChunkSize)
	ch, ok := w.chunks[key]
	if !ok {
		panic("expected look at voxel to be in a loaded chunk")
	}
	vox := ch.GetVoxelFromFlatData(v.Pos)
	pos := vox.Pos.AsVec3()
	w.selectedVoxel.SetData([]float32{pos.X(), pos.Y(), pos.Z(), float32(vox.GetVbits())})
}

// SetVoxel updates a voxel's variables in the world if the chunk
// that it would belong to is currently loaded.
func (w *World) SetVoxel(v *Voxel) {
	key := v.Pos.GetChunkPos(ChunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	chunk, ok := w.chunks[key]
	if !ok {
		return
	}
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
		ch, ok := w.chunks[k]
		if !ok {
			w.pendingFuncs[k] = append(w.pendingFuncs[k], adjModFn(mod))
		} else {
			ch.AddAdjacency(mod.off, mod.adjDiff)
		}
	}
	// TODO do this in a goroutine
	// w.updateLighting()
}

func (w *World) RemoveVoxel(v VoxelPos) {
	key := v.GetChunkPos(ChunkSize)
	chunk, ok := w.chunks[key]
	if !ok {
		return
	}
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
		ch, ok := w.chunks[k]
		if !ok {
			w.pendingFuncs[k] = append(w.pendingFuncs[k], adjModFn(mod))
		} else {
			ch.RemoveAdjacency(mod.off, mod.adjDiff)
		}
	}
	// TODO do this in a goroutine
	// w.updateLighting()
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

func (w *World) updateLighting() {
	for _, ch := range w.chunks {
		for i := 0; i < ch.size; i++ {
			for j := 0; j < ch.size; j++ {
				for k := 0; k < ch.size; k++ {
					p := VoxelPos{
						X: ch.AsVoxelPos().X + i,
						Y: ch.AsVoxelPos().Y + j,
						Z: ch.AsVoxelPos().Z + k,
					}
					v := ch.GetVoxelFromFlatData(p)
					v.LightBits = 0
					ch.SetVoxel(&v)
				}
			}
		}
	}
	for _, ch := range w.chunks {
		for lightPos := range ch.lights {
			v := ch.GetVoxelFromFlatData(lightPos)
			v.SetLightValue(MaxLightValue, LightLeft)
			v.SetLightValue(MaxLightValue, LightRight)
			v.SetLightValue(MaxLightValue, LightBack)
			v.SetLightValue(MaxLightValue, LightFront)
			v.SetLightValue(MaxLightValue, LightTop)
			v.SetLightValue(MaxLightValue, LightBottom)
			ch.SetVoxel(&v)

			w.lightFrom(ch, lightPos, MaxLightValue)
			// TODO avoid resetting by storing identifier for light block to know if it owns the explored bit
			// identifier could be local coordinate
			// log.Debug(ch)
		}
	}
	for _, ch := range w.chunks {
		for i := 0; i < ch.size; i++ {
			for j := 0; j < ch.size; j++ {
				for k := 0; k < ch.size; k++ {
					p := VoxelPos{
						X: ch.AsVoxelPos().X + i,
						Y: ch.AsVoxelPos().Y + j,
						Z: ch.AsVoxelPos().Z + k,
					}
					v := ch.GetVoxelFromFlatData(p)
					v.SetLightValue(0, LightValue)
					ch.SetVoxel(&v)
				}
			}
		}
	}
}

func (w *World) lightFrom(c *Chunk, p VoxelPos, value uint32) {
	if value < 0 || value > MaxLightValue || !c.IsWithinChunk(p) {
		panic("improper usage: bad lighting value or p outside chunk")
	}
	if value == 0 {
		return
	}
	currBlock := c.GetVoxelFromFlatData(p)
	if currBlock.GetLightValue(LightValue) >= value {
		return
	}
	if currBlock.Btype == Air {
		currBlock.SetLightValue(value, LightValue)
		c.SetVoxel(&currBlock)
	}
	dirs := []struct {
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
	for _, dir := range dirs {
		offP := p.Add(dir.off)
		offKey := offP.GetChunkPos(ChunkSize)
		if offCh, ok := w.chunks[offKey]; ok {
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			if adjBlock.Btype == Air {
				// continue search down open path
				w.lightFrom(offCh, offP, value-1)
			} else {
				// path is blocked off, apply lighting value to the face
				if adjBlock.GetLightValue(dir.face) < value {
					// log.Debugf("Setting face %v of %v to value %v", GetLightMaskName(dir.face), offP, value)
					adjBlock.SetLightValue(value, dir.face)
					offCh.SetVoxel(&adjBlock)
				}
			}
		} else {
			// do nothing intentionally
		}
	}
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
				sw := util.Start()
				objs, err := voxgl.NewChunkObject(nil)
				if err != nil {
					panic(fmt.Sprint(err))
				}
				ch.SetObjs(objs)
				w.chunks[ch.Pos] = ch
				w.applyPendingMods(ch)

				// TODO do this in a goroutine
				// w.updateLighting()
				sw.StopRecordAverage("Chunk objs + pending mods")
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
	_, loaded := w.chunks[key]
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
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)
	for {
		select {
		case key := <-w.saved:
			delete(w.chunkSaving, key)
			if rng.Contains(key) {
				w.chunkExpect[key] = struct{}{}
				w.requestChunk(key)
			}
		default:
			return
		}
	}
}

func (w *World) applyPendingMods(ch *Chunk) {
	if fns, hasPending := w.pendingFuncs[ch.Pos]; hasPending {
		for _, fn := range fns {
			fn(ch)
		}
		delete(w.pendingFuncs, ch.Pos)
	}
}

func adjModFn(mod AdjMod) func(*Chunk) {
	return func(ch *Chunk) {
		if mod.add {
			ch.AddAdjacency(mod.off, mod.adjDiff)
		} else {
			ch.RemoveAdjacency(mod.off, mod.adjDiff)
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
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)
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
	for key, ch := range w.chunks {
		if _, ok := w.chunkExpect[key]; !ok {
			w.chunks[key].Destroy()
			delete(w.chunks, key)
			if ch.modified {
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
				}(key, ch)
			}
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

	for _, chunk := range w.chunks {
		w.cubeMap.Bind()
		chunk.Render(w.cam)
		w.cubeMap.Unbind()
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
	// TODO make a more well defined cleanup routine
	w.cancel = true
	w.cacheLock.Lock()
	for key := range w.pendingFuncs {
		if ch, ok := w.chunks[key]; !ok {
			chunk, loaded := w.cache.Load(key)
			if !loaded {
				chunk = NewChunk(ChunkSize, key, w.gen)
			}
			w.applyPendingMods(chunk)
			w.cache.Save(chunk)
		} else {
			w.applyPendingMods(ch)
			w.cache.Save(ch)
			delete(w.chunks, key)
		}
	}
	for _, ch := range w.chunks {
		if ch.modified {
			w.cache.Save(ch)
		}
	}
	w.cache.Destroy()
}

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
	w.saveRoutine()
	w.cacheLock.Unlock()
}
