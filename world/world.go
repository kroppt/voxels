package world

import (
	"fmt"
	"math/rand"
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
	ubo          *gfx.BufferObject
	cam          *Camera
	chunks       map[ChunkPos]*Chunk
	chunkExpect  map[ChunkPos]struct{}
	chunkSaving  map[ChunkPos]struct{}
	chunkLoading map[ChunkPos]struct{}
	chunkChan    chan *Chunk
	saved        chan ChunkPos
	loaded       chan ChunkPos
	cubeMap      *gfx.CubeMap
	gen          Generator
	cache        *Cache
}

const ChunkSize = 3
const chunkRenderDist = 5
const cacheThreshold = 10

// New returns a new world.World.
func New() *World {
	ubo := gfx.NewBufferObject()
	var mat glm.Mat4
	// opengl memory allocation, 2x mat4 = 1 for proj + 1 for view
	ubo.BufferData(gl.UNIFORM_BUFFER, uint32(2*unsafe.Sizeof(mat)), gl.Ptr(nil), gl.STATIC_DRAW)
	// use binding = 0
	ubo.BindBufferBase(gl.UNIFORM_BUFFER, 0)
	cam := NewCameraDefault()

	world := &World{
		ubo:       ubo,
		cam:       cam,
		chunkChan: make(chan *Chunk),
		saved:     make(chan ChunkPos),
		loaded:    make(chan ChunkPos),
		gen:       FlatWorldGenerator{},
	}

	cam.SetPosition(&glm.Vec3{0.5, 7.5, 2})
	// cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})

	cache, err := NewCache("world_meta", "world_data", cacheThreshold)
	if err != nil {
		panic(fmt.Sprint(err))
	}
	world.cache = cache

	rand.Seed(time.Now().UnixNano())
	chunks := make(map[ChunkPos]*Chunk)
	world.chunkExpect = make(map[ChunkPos]struct{})
	world.chunkLoading = make(map[ChunkPos]struct{})

	world.chunks = chunks
	world.cubeMap = loadSpriteSheet("sprite_sheet.png")
	world.chunkSaving = make(map[ChunkPos]struct{})

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

// SetVoxel updates a voxel's variables in the world if the chunk
// that it would belong to is currently loaded.
func (w *World) SetVoxel(v *Voxel) {
	key := v.Pos.GetChunkPos(ChunkSize)
	// log.Debugf("Adding voxel at %v in chunk %v", v.Pos, key)
	chunk, ok := w.chunks[key]
	if !ok {
		return
	}
	target := chunk.GetVoxel(v.Pos)
	v.AdjMask = target.AdjMask
	chunk.SetVoxel(v)
	chunk.SetModified()
	{
		p := VoxelPos{v.Pos.X - 1, v.Pos.Y, v.Pos.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentRight)
	}
	{
		p := VoxelPos{v.Pos.X + 1, v.Pos.Y, v.Pos.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentLeft)
	}
	{
		p := VoxelPos{v.Pos.X, v.Pos.Y - 1, v.Pos.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentTop)
	}
	{
		p := VoxelPos{v.Pos.X, v.Pos.Y + 1, v.Pos.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentBottom)
	}
	{
		p := VoxelPos{v.Pos.X, v.Pos.Y, v.Pos.Z - 1}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentBack)
	}
	{
		p := VoxelPos{v.Pos.X, v.Pos.Y, v.Pos.Z + 1}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.AddAdjacency(p, AdjacentFront)
	}
}

func (w *World) RemoveVoxel(v VoxelPos) {
	key := v.GetChunkPos(ChunkSize)
	chunk, ok := w.chunks[key]
	if !ok {
		return
	}
	// log.Debugf("removing voxel in chunk %v", key)
	chunk.SetVoxel(&Voxel{
		Pos:   v,
		Btype: Air,
	})
	chunk.SetModified()
	chunk.root, _ = chunk.root.Remove(v)
	{
		p := VoxelPos{v.X - 1, v.Y, v.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentRight)
	}
	{
		p := VoxelPos{v.X + 1, v.Y, v.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentLeft)
	}
	{
		p := VoxelPos{v.X, v.Y - 1, v.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentTop)
	}
	{
		p := VoxelPos{v.X, v.Y + 1, v.Z}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentBottom)
	}
	{
		p := VoxelPos{v.X, v.Y, v.Z - 1}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentBack)
	}
	{
		p := VoxelPos{v.X, v.Y, v.Z + 1}
		k := p.GetChunkPos(ChunkSize)
		ch, ok := w.chunks[k]
		if !ok {
			panic("the player unlocked TNT and wiremod")
		}
		ch.RemoveAdjacency(p, AdjacentFront)
	}
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

func (w *World) receiveExpectedAsync() {
	for {
		select {
		case ch := <-w.chunkChan:
			if _, ok := w.chunkExpect[ch.Pos]; ok {
				// if _, loaded := w.chunks[ch.Pos]; loaded {
				// 	continue
				// }
				// the chunk has arrived and we expected it
				// give the chunk its object
				sw := util.Start()
				objs, err := voxgl.NewColoredObject(nil)
				if err != nil {
					panic(fmt.Sprint(err))
				}
				ch.SetObjs(objs)
				w.chunks[ch.Pos] = ch
				sw.StopRecordAverage("Chunk objs")
			}
		default:
			return
		}
	}
}

func (w *World) requestExpectedChunks() {
	for key := range w.chunkExpect {
		if _, loaded := w.chunks[key]; loaded {
			continue
		}
		if _, loading := w.chunkLoading[key]; loading {
			continue
		}
		w.chunkLoading[key] = struct{}{}
		go func(key ChunkPos) {
			// check cache for a saved chunk
			chunk, loaded := w.cache.Load(key)
			if !loaded {
				chunk = NewChunk(ChunkSize, key, w.gen)
			}
			// TODO switching these lines changes things
			// should these two be tied?
			w.chunkChan <- chunk
			w.loaded <- key
		}(key)
	}
}

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

func (w *World) evictUnexpectedChunks() {
	for key, ch := range w.chunks {
		if _, ok := w.chunkExpect[key]; !ok {
			if _, loaded := w.chunks[key]; !loaded {
				// TODO is this possible?
				panic("not expected and not loaded")
				continue
			}
			w.chunks[key].Destroy()
			delete(w.chunks, key)
			if ch.modified {
				w.chunkSaving[key] = struct{}{}
				go func(key ChunkPos, ch *Chunk) {
					w.cache.Save(ch)
					w.saved <- key
				}(key, ch)
			}
		}
	}
}

// Render renders the chunks of the world in OpenGL.
func (w *World) Render() error {
	w.evictUnexpectedChunks()
	w.checkSavingStatus()
	w.updateExpectedChunks()
	w.checkLoadingStatus()
	w.requestExpectedChunks()
	w.receiveExpectedAsync()

	if w.cam.IsDirty() {
		err := w.updateUBO()
		if err != nil {
			return err
		}
		w.cam.Clean()
	}

	culled := 0
	for _, chunk := range w.chunks {
		if w.cam.IsWithinFrustum(chunk.AsVoxelPos().AsVec3(), float32(chunk.size), float32(chunk.size), float32(chunk.size)) {
			w.cubeMap.Bind()
			chunk.Render()
			w.cubeMap.Unbind()
		} else {
			culled++
		}
	}
	// log.Debugf("culled %v / %v chunks", culled, len(w.chunks))
	return nil
}

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
	w.cache.Destroy()
}
