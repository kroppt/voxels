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
	chunksLoaded  map[ChunkPos]*LoadedChunk
	chunksSaving  map[ChunkPos]struct{}
	chunksLoading map[ChunkPos]struct{}
	currChunk     ChunkPos
	chunkChan     chan *Chunk
	processed     chan ChunkPos
	saved         chan ChunkPos
	cubeMap       *gfx.CubeMap
	gen           Generator
	cache         *Cache
	cacheLock     sync.Mutex
	chunkLock     sync.RWMutex

	// cancel        bool
	selectedVoxel *voxgl.Object
	selected      bool
	crosshair     *voxgl.Crosshair
}

// LoadedChunk contains a loaded chunk and booleans
// describing various states that the chunk is in, and
// operations that are expected to happen on it.
type LoadedChunk struct {
	chunk          *Chunk
	modified       bool
	doRender       bool
	processing     bool
	relightActions []func()
}

const ChunkSize = 8
const chunkRenderDist = 4

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
		// loaded:    make(chan ChunkPos),
		processed: make(chan ChunkPos),
		gen:       AlexGenerator{},
		crosshair: crosshair,
		cacheLock: sync.Mutex{},
		chunkLock: sync.RWMutex{},
	}

	cam.SetPosition(&glm.Vec3{0.5, 30.5, 0.5})
	cam.LookAt(&glm.Vec3{0.5, 0.5, 2})

	cache, err := NewCache("world_meta", "world_data", cacheThreshold)
	if err != nil {
		panic(fmt.Sprint(err))
	}
	world.cache = cache

	rand.Seed(time.Now().UnixNano())

	world.chunksLoaded = map[ChunkPos]*LoadedChunk{}
	world.chunksSaving = make(map[ChunkPos]struct{})
	world.chunksLoading = make(map[ChunkPos]struct{})

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
	w.chunkLock.RLock()
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
	w.chunkLock.RUnlock()
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
	w.chunkLock.RLock()
	loadedCh, ok := w.chunksLoaded[key]
	w.chunkLock.RUnlock()
	if !ok {
		panic("expected look at voxel to be in a loaded chunk")
	}
	vox := loadedCh.chunk.GetVoxelFromFlatData(v.Pos)
	pos := vox.Pos.AsVec3()
	w.selectedVoxel.SetData([]float32{pos.X(), pos.Y(), pos.Z(), float32(vox.GetVbits())})
}

func (w *World) getUniqueSources(ch *Chunk, p VoxelPos) map[VoxelPos]uint32 {
	srcMap := ch.lightRefs[p]
	uniques := make(map[VoxelPos]uint32)
	for src, val := range srcMap {
		uniques[src] = val
	}
	for _, mod := range faceMods {
		offP := p.Add(mod.off)
		offKey := offP.GetChunkPos(ChunkSize)
		w.chunkLock.RLock()
		if loadedCh, ok := w.chunksLoaded[offKey]; ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			// TODO remove these checks, look into all blocks? (how could a solid )
			if adjBlock.Btype == Air || adjBlock.Btype == Light {
				offSrcMap := offCh.lightRefs[offP]
				for offSrc, offVal := range offSrcMap {
					uniques[offSrc] = offVal
				}
			}
		}
		w.chunkLock.RUnlock()
	}
	return uniques
}

func (w *World) blockPlaceFn(ch *Chunk, p VoxelPos, light bool) func() {
	/*
		Block is placed (could block light)
			- collect all unique light sources from 6 neighbors
			- lightRemoveFrom on every unique light
			- lightFrom on every unique light
			if block was light source:
				- lightFrom() <- // TODO caught in uniques ?
	*/
	return func() {
		uniques := w.getUniqueSources(ch, p)
		for uniqueSrc := range uniques {
			srcChKey := uniqueSrc.GetChunkPos(ChunkSize)
			w.chunkLock.RLock()
			srcCh := w.chunksLoaded[srcChKey].chunk
			w.chunkLock.RUnlock()
			// TODO have custom light value per light so its not hacked in here
			w.lightRemoveFrom(srcCh, uniqueSrc, uniqueSrc)
			w.lightFrom(srcCh, uniqueSrc, MaxLightValue, uniqueSrc)
		}
		if light {
			w.lightFrom(ch, p, MaxLightValue, p)
		}
	}
}

func (w *World) blockRemoveFn(ch *Chunk, p VoxelPos, light bool) func() {
	/*
		Block is removed (could allow light to pass through)
		if block was light source:
			- lightRemoveFrom()
		- relight from all unique light sources from 6 neighbors

	*/
	return func() {
		if light {
			w.lightRemoveFrom(ch, p, p)
		}
		uniques := w.getUniqueSources(ch, p)
		for uniqueSrc := range uniques {
			srcChKey := uniqueSrc.GetChunkPos(ChunkSize)
			w.chunkLock.RLock()
			srcCh := w.chunksLoaded[srcChKey].chunk
			w.chunkLock.RUnlock()
			// TODO have custom light value per light so its not hacked in here
			w.lightFrom(srcCh, uniqueSrc, MaxLightValue, uniqueSrc)
		}
	}
}

func (w *World) relightAllFn(ch *Chunk) func() {
	/*
		Chunk is loaded
			- lightFrom() on every light source in chunk

		Chunk transitions from buffered(=in invisible buffer range) to rendered
			- lightFrom() on every light source in chunk
	*/
	return func() {
		var lightCopy []VoxelPos
		ch.lightLock.RLock()
		for lightPos := range ch.lights {
			lightCopy = append(lightCopy, lightPos)
		}
		ch.lightLock.RUnlock()
		for _, light := range lightCopy {
			w.lightFrom(ch, light, MaxLightValue, light)
		}
	}
}

// SetVoxel updates a voxel's variables in the world
func (w *World) SetVoxel(v *Voxel) {
	// TODO setting a voxel to air is *similar* to removing, but not the same
	// -> what would setting a voxel to air do here? should this be prevented?
	key := v.Pos.GetChunkPos(ChunkSize)
	if !w.hasSurroundingChunks(key) {
		return
	}
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	w.chunkLock.RLock()
	loadedCh, ok := w.chunksLoaded[key]
	w.chunkLock.RUnlock()
	if !ok {
		panic(fmt.Sprintf("tried to set a voxel (%v) that belongs to an unrendered chunk", v.Pos))
	}
	chunk := loadedCh.chunk
	target := chunk.GetVoxelFromFlatData(v.Pos)
	v.AdjMask = target.AdjMask
	chunk.SetVoxel(v)
	for _, mod := range faceMods {
		mod.off = mod.off.Add(v.Pos)
		k := mod.off.GetChunkPos(ChunkSize)
		w.chunkLock.RLock()
		ch, ok := w.chunksLoaded[k]
		w.chunkLock.RUnlock()
		if !ok {
			// TODO how did this happen
			panic("not currently handling pending mods to unloaded chunks")
		} else {
			ch.chunk.AddAdjacency(mod.off, mod.adjFace)
			ch.modified = true
		}
	}

	isLight := chunk.GetVoxelFromFlatData(v.Pos).Btype == Light
	loadedCh.relightActions = append(loadedCh.relightActions, w.blockPlaceFn(chunk, v.Pos, isLight))
}

// RemoveVoxel sets a voxel to air and updates necessary structures.
// - Remove from the octree of its parent chunk so it's not considered in intersections
// - Update adjacency bits of surrounding voxels so you can see their sides
// - Trigger a relight
func (w *World) RemoveVoxel(v VoxelPos) {
	key := v.GetChunkPos(ChunkSize)
	if !w.hasSurroundingChunks(key) {
		return
	}
	w.chunkLock.RLock()
	loadedCh, ok := w.chunksLoaded[key]
	w.chunkLock.RUnlock()
	if !ok {
		panic(fmt.Sprintf("tried to remove a voxel (%v) that belongs to an unloaded chunk", v))
	}
	chunk := loadedCh.chunk
	// TODO setting voxel to air and removing node from octree are tied
	// together in a "delete voxel" operation, this should be more well defined
	// inside chunk.go in a single function call
	isLight := chunk.GetVoxelFromFlatData(v).Btype == Light
	chunk.SetVoxel(&Voxel{
		Pos:     v,
		Btype:   Air,
		AdjMask: AdjacentNone,
	})
	chunk.root, _ = chunk.root.Remove(v)
	for _, mod := range faceMods {
		mod.off = mod.off.Add(v)
		k := mod.off.GetChunkPos(ChunkSize)
		w.chunkLock.RLock()
		loadedCh, ok := w.chunksLoaded[k]
		w.chunkLock.RUnlock()
		if !ok {
			panic("not currently handling pending mods to unloaded chunks")
		} else {
			loadedCh.chunk.RemoveAdjacency(mod.off, mod.adjFace)
			loadedCh.modified = true
		}
	}

	loadedCh.relightActions = append(loadedCh.relightActions, w.blockRemoveFn(chunk, v, isLight))
}

// GetCamera returns a reference to the camera.
func (w *World) GetCamera() *Camera {
	return w.cam
}

// TODO Create UBO extraction
func (w *World) updateUBO() {
	cam := *w.GetCamera()
	view := cam.GetViewMat()
	err := w.ubo.BufferSubData(gl.UNIFORM_BUFFER, 0, uint32(unsafe.Sizeof(view)), gl.Ptr(&view[0]))
	if err != nil {
		panic(fmt.Sprintf("failed to upload camera view to ubo: %v", err))
	}
	proj := cam.GetProjMat()
	err = w.ubo.BufferSubData(gl.UNIFORM_BUFFER, uint32(unsafe.Sizeof(view)), uint32(unsafe.Sizeof(proj)), gl.Ptr(&proj[0]))
	if err != nil {
		panic(fmt.Sprintf("failed to upload camera proj to ubo: %v", err))
	}
}

func (w *World) requestAsyncLighting() {
	w.chunkLock.RLock()
	// TODO instead of looping over everything that is loaded, maintain a list of chunks
	// that have relight requests, because only those need to be considered
	for key, loadedCh := range w.chunksLoaded {
		if len(loadedCh.relightActions) > 0 && w.hasSurroundingChunks(key) && !w.hasProcessingNeighbor(key) {
			w.setNeighborsProcessing(key, true)
			go func(key ChunkPos, tasks []func()) {
				for _, task := range tasks {
					task()
				}
				w.processed <- key
			}(loadedCh.chunk.Pos, loadedCh.relightActions)
			loadedCh.relightActions = nil
		}
	}
	w.chunkLock.RUnlock()
}

var faceMods = [6]struct {
	off       VoxelPos
	lightFace LightMask
	adjFace   AdjacentMask
}{
	{VoxelPos{-1, 0, 0}, LightRight, AdjacentRight},
	{VoxelPos{1, 0, 0}, LightLeft, AdjacentLeft},
	{VoxelPos{0, -1, 0}, LightTop, AdjacentTop},
	{VoxelPos{0, 1, 0}, LightBottom, AdjacentBottom},
	{VoxelPos{0, 0, -1}, LightBack, AdjacentBack},
	{VoxelPos{0, 0, 1}, LightFront, AdjacentFront},
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
	for _, mod := range faceMods {
		offP := p.Add(mod.off)
		offKey := offP.GetChunkPos(ChunkSize)
		w.chunkLock.RLock()
		loadedCh, ok := w.chunksLoaded[offKey]
		w.chunkLock.RUnlock()
		if ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			_, secondLargest, found := c.GetBrightestSource(p)
			if adjBlock.Btype != Light {
				if found {
					adjBlock.SetLightValue(secondLargest, mod.lightFace)
				} else {
					adjBlock.SetLightValue(0, mod.lightFace)
				}
			} else {
				// if it is a light, for now just don't do anything
				// maybe in the future, dim lights could have their sides
				// lit up more by a different, stronger light source, more than
				// they would otherwise light themselves
			}
			offCh.SetVoxelLightBits(adjBlock)
			w.lightRemoveFrom(offCh, offP, src)
		}
	}
}

// lightFrom peforms a flood fill recursive algorithm that propagates light from the
// specified VoxelPos within initial light intensity as "value" coming from light
// source "src". The light decreases in intensity as the traversal propagates further
// from the light source, applying its value to the faces of voxels that it passes by.
// TODO try BFS instead of DFS
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
		v := Voxel{
			Pos: p,
		}
		v.SetLightValue(MaxLightValue, LightLeft)
		v.SetLightValue(MaxLightValue, LightRight)
		v.SetLightValue(MaxLightValue, LightBack)
		v.SetLightValue(MaxLightValue, LightFront)
		v.SetLightValue(MaxLightValue, LightTop)
		v.SetLightValue(MaxLightValue, LightBottom)
		c.SetVoxelLightBits(v)
	}
	for _, mod := range faceMods {
		offP := p.Add(mod.off)
		offKey := offP.GetChunkPos(ChunkSize)
		w.chunkLock.RLock()
		loadedCh, ok := w.chunksLoaded[offKey]
		w.chunkLock.RUnlock()
		if ok {
			offCh := loadedCh.chunk
			adjBlock := offCh.GetVoxelFromFlatData(offP)
			// TODO transparency abstraction, opacity for value reduction other than -1
			if adjBlock.Btype == Air {
				// continue search down open path
				w.lightFrom(offCh, offP, value-1, src)
			} else {
				// path is blocked off, apply lighting value to the face *if brighter*
				if adjBlock.GetLightValue(mod.lightFace) < value {
					adjBlock.SetLightValue(value, mod.lightFace)
					offCh.SetVoxelLightBits(adjBlock)
				}
			}
		}
	}
}

// isWithinRenderDist returns whether the key ChunkPos is within render distance
// from the chunk that the camera is currently in.
func (w *World) isWithinRenderDist(key ChunkPos) bool {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)
	return rng.Contains(key)
}

// hasSurroundingChunks returns whether the all "surrounding" chunks around
// the chunk specified by the key have loaded. "Surrounding" means all chunks within
// the buffer distance around the chunk (chunkRenderBuffer), including the key chunk itself,
// with the exception being chunks past invisible "buffer" chunks, which are excluded
// from the check and not required to be loaded for invisible chunks to consider themselves surrounded.
func (w *World) hasSurroundingChunks(key ChunkPos) bool {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	// loadedRng handles the edge case for buffered, invisible
	// chunks that aren't expected to have chunks beyond them
	// loaded before they run their lighting algorithm.
	loadRange := currChunk.GetSurroundings(chunkRenderDist + chunkRenderBuffer)
	localBufRng := key.GetSurroundings(chunkRenderBuffer)
	allHere := true
	localBufRng.ForEach(func(pos ChunkPos) {
		if !loadRange.Contains(pos) {
			return
		}
		w.chunkLock.RLock()
		if _, loaded := w.chunksLoaded[pos]; !loaded {
			allHere = false
		}
		w.chunkLock.RUnlock()
	})
	return allHere
}

// hasProcessingNeighbor returns whether any loaded "surrounding" chunk (see above)
// is still processing, including the key chunk itself.
func (w *World) hasProcessingNeighbor(key ChunkPos) bool {
	localBufRng := key.GetSurroundings(chunkRenderBuffer)
	processing := false
	localBufRng.ForEach(func(pos ChunkPos) {
		w.chunkLock.RLock()
		if ch, loaded := w.chunksLoaded[pos]; loaded && ch.processing {
			processing = true
		}
		w.chunkLock.RUnlock()
	})
	return processing
}

// setNeighborsProcessing sets the processing state to the specified value
// for every loaded "surrounding" chunk (see above), including the key chunk itself.
func (w *World) setNeighborsProcessing(key ChunkPos, processing bool) {
	localBufRng := key.GetSurroundings(chunkRenderBuffer)
	localBufRng.ForEach(func(pos ChunkPos) {
		w.chunkLock.RLock()
		if loadedCh, loaded := w.chunksLoaded[pos]; loaded {
			loadedCh.processing = processing
		}
		w.chunkLock.RUnlock()
	})
}

// setNeighborsDirty sets the dirty bit to the specified value
// for every loaded "surrounding" chunk (see above), including the key chunk itself.
func (w *World) setNeighborsDirty(key ChunkPos, dirty bool) {
	localBufRng := key.GetSurroundings(chunkRenderBuffer)
	localBufRng.ForEach(func(pos ChunkPos) {
		w.chunkLock.RLock()
		if loadedCh, loaded := w.chunksLoaded[pos]; loaded {
			loadedCh.chunk.dirty = dirty
		}
		w.chunkLock.RUnlock()
	})
}

// receiveAllChannels is only called by the main thread, reading off all channels, including:
// - saved: indicates that a chunk finished saving by sending its key
// - processed: indicates that a chunk finished processing by sending its key
//		- (right now, the only async processing is lighting calculations)
// - chunkChan: indicates that a chunk finished loading by sending the chunk
func (w *World) receiveAllChannels() {
	for {
		select {
		case key := <-w.saved:
			delete(w.chunksSaving, key)
		case key := <-w.processed:
			w.setNeighborsProcessing(key, false)
			w.setNeighborsDirty(key, true)
		case ch := <-w.chunkChan:
			delete(w.chunksLoading, ch.Pos)
			if w.isExpected(ch.Pos) {
				// the chunk has arrived and we expected it
				// give the chunk its object
				objs, err := voxgl.NewChunkObject(nil)
				if err != nil {
					panic(fmt.Sprintf("failure making NewChunkObject: %v", err))
				}
				ch.SetObjs(objs)

				loadedCh := LoadedChunk{
					chunk:          ch,
					modified:       false,
					processing:     false,
					doRender:       w.isWithinRenderDist(ch.Pos) && w.hasSurroundingChunks(ch.Pos),
					relightActions: []func(){w.relightAllFn(ch)},
				}

				w.chunkLock.Lock()
				w.chunksLoaded[ch.Pos] = &loadedCh
				w.chunkLock.Unlock()
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
// chunk channel.
// A request to any chunk that is already loaded, loading, or saving will
// be ignored.
func (w *World) requestChunk(key ChunkPos) {
	w.chunkLock.RLock()
	_, loaded := w.chunksLoaded[key]
	w.chunkLock.RUnlock()
	_, loading := w.chunksLoading[key]
	_, saving := w.chunksSaving[key]
	if loaded || loading || saving {
		return
	}
	w.chunksLoading[key] = struct{}{}
	go func(key ChunkPos) {
		// check cache for a saved chunk
		w.cacheLock.Lock()
		// if w.cancel {
		// 	w.cacheLock.Unlock()
		// 	return
		// }
		chunk, loaded := w.cache.Load(key)
		w.cacheLock.Unlock()
		if !loaded {
			chunk = NewChunk(ChunkSize, key, w.gen)
		}
		w.chunkChan <- chunk
	}(key)
}

func (w *World) isExpected(key ChunkPos) bool {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist + chunkRenderBuffer)
	return rng.Contains(key)
}

// requestExpectedChunks attempts to request every expected chunk
// Note that in requestChunk, it will ignore requests to loaded,
// loading or saving chunks.
func (w *World) requestExpectedChunks() {
	// slightly larger range for chunks that should be loaded
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist + chunkRenderBuffer)
	rng.ForEach(func(pos ChunkPos) {
		w.requestChunk(pos)
	})
}

// updateLoadedChunks checks all loaded chunks that are being rendered
// and removes any that are no longer expected. If they were modified, it
// puts in an asynchronous request for the chunk to be saved, marking its key
// as in the process of saving, and indicating completion on the saved channel
// when it is done.
func (w *World) updateLoadedChunks() {
	w.chunkLock.Lock()
	for key, loadedCh := range w.chunksLoaded {
		if !w.isExpected(key) && !loadedCh.processing {
			w.chunksLoaded[key].chunk.Destroy()
			delete(w.chunksLoaded, key)
			if loadedCh.modified {
				w.chunksSaving[key] = struct{}{}
				go func(key ChunkPos, ch *Chunk) {
					w.cacheLock.Lock()
					// if w.cancel {
					// 	w.cacheLock.Unlock()
					// 	return
					// }
					w.cache.Save(ch)
					w.cacheLock.Unlock()
					w.saved <- key
				}(key, loadedCh.chunk)
			}
		} else {
			shouldRender := w.isWithinRenderDist(key)
			if !loadedCh.doRender && shouldRender {
				// this chunk was a buffer chunk and has moved into render dist
				loadedCh.relightActions = append(loadedCh.relightActions, w.relightAllFn(loadedCh.chunk))
			}
			loadedCh.doRender = shouldRender
		}
	}
	w.chunkLock.Unlock()
}

func (w *World) update() {
	w.updateLoadedChunks()
	w.requestExpectedChunks()
}

// Render renders the chunks of the world in OpenGL.
func (w *World) Render() error {
	sw := util.Start()
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	if currChunk != w.currChunk {
		w.currChunk = currChunk
		// TODO use w.update() in here
	}
	// TODO figure out how to fix no update until u move chunks if update is above
	w.update()

	w.requestAsyncLighting()
	w.receiveAllChannels()
	sw.StopRecordAverage("total update logic")

	if w.cam.IsDirty() {
		w.updateUBO()
		w.cam.Clean()
	}

	w.chunkLock.RLock()
	w.cubeMap.Bind()
	for _, loadedCh := range w.chunksLoaded {
		if loadedCh.doRender {
			loadedCh.chunk.Render(w.cam)
		}
	}
	w.cubeMap.Unbind()
	w.chunkLock.RUnlock()

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
	// TODO cancel properly... lock around this too?
	// w.cancel = true
	w.cacheLock.Lock()
	w.chunkLock.RLock()
	for _, ch := range w.chunksLoaded {
		if ch.modified {
			w.cache.Save(ch.chunk)
		}
	}
	w.chunkLock.RUnlock()
	w.cache.Sync()
	w.cache.Destroy()
	w.cacheLock.Unlock()
}

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
	w.saveRoutine()
}
