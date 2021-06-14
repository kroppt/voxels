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
	ubo         *gfx.BufferObject
	cam         *Camera
	chunks      map[ChunkPos]*Chunk
	chunkExpect map[ChunkPos]struct{}
	lastChunk   ChunkPos
	chunkChan   chan *Chunk
	cubeMap     *gfx.CubeMap
	gen         Generator
	cache       *Cache
}

const ChunkSize = 2
const chunkRenderDist = 5

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
		gen:       FlatWorldGenerator{},
	}

	cam.SetPosition(&glm.Vec3{0.5, 7.5, 2})
	// cam.LookAt(&glm.Vec3{0.5, 0.5, 0.5})
	currChunk := cam.AsVoxelPos().GetChunkPos(ChunkSize)
	rng := currChunk.GetSurroundings(chunkRenderDist)

	cache, err := NewCache("world_meta", "world_data")
	if err != nil {
		panic(fmt.Sprint(err))
	}
	world.cache = cache

	rand.Seed(time.Now().UnixNano())
	chunks := make(map[ChunkPos]*Chunk)
	world.chunkExpect = make(map[ChunkPos]struct{})
	rng.ForEach(func(pos ChunkPos) {
		world.chunkExpect[pos] = struct{}{}
		world.LoadChunkAsync(pos)
	})

	world.chunks = chunks
	world.lastChunk = currChunk
	world.cubeMap = loadSpriteSheet("sprite_sheet.png")

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
	// log.Debug(len(sprytes))
	// log.Debugf("w = %v, h = %v, layers = %v", w, h, layers)
	texAtlas, err := gfx.NewCubeMap(w, layers, sprytes, gl.RGBA, 4, 4)
	if err != nil {
		panic("failed to create 3d texture")
	}
	texAtlas.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	texAtlas.SetParameter(gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	return &texAtlas
}

// LoadChunkAsync starts a thread that will put the chunk in the chunkChan
func (w *World) LoadChunkAsync(pos ChunkPos) {
	// immediately set that the chunk is expected to be loaded
	w.chunkExpect[pos] = struct{}{}
	go func() {
		// check cache for a saved chunk
		chunk, loaded := w.cache.Load(pos)
		if !loaded {
			chunk = NewChunk(ChunkSize, pos, w.gen)
		}

		w.chunkChan <- chunk
	}()
}

func (w *World) UpdateChunksAsync() {
	for {
		select {
		case ch := <-w.chunkChan:
			if _, ok := w.chunkExpect[ch.Pos]; ok {
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
			// the load was too late, chunk already unloaded
			// before it was even loaded for the first time
			// do not destroy, because it was never created

		default:
			return
		}
	}
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

// Destroy frees external resources.
func (w *World) Destroy() {
	w.ubo.Destroy()
	w.cache.Destroy()
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

func (w *World) UpdateChunks() {
	currChunk := w.cam.AsVoxelPos().GetChunkPos(ChunkSize)
	if currChunk == w.lastChunk {
		return
	}
	// the camera position has moved chunks
	// load new chunks
	rng := currChunk.GetSurroundings(chunkRenderDist)
	counter := 0
	rng.ForEach(func(pos ChunkPos) {
		if _, ok := w.chunkExpect[pos]; !ok {
			// we do not yet expect chunk i,j to be loaded,
			// but we want to add it as a new one
			w.LoadChunkAsync(pos)
			counter++
		}
	})
	// delete old chunks
	lastRng := w.lastChunk.GetSurroundings(chunkRenderDist)
	lastRng.ForEach(func(pos ChunkPos) {
		inOld := lastRng.Contains(pos)
		inNew := rng.Contains(pos)
		if inOld && !inNew {
			if ch, ok := w.chunks[pos]; ok {
				// save chunk if modified
				if ch.modified {
					w.cache.Save(ch)
				}
				// the chunk-to-be-unloaded actually exists
				w.chunks[pos].Destroy()
				delete(w.chunks, pos)
				delete(w.chunkExpect, pos)
			}
		}
	})
	w.lastChunk = currChunk
}

// Render renders the chunks of the world in OpenGL.
func (w *World) Render() error {
	w.UpdateChunksAsync()

	if w.cam.IsDirty() {
		err := w.updateUBO()
		if err != nil {
			return err
		}
		w.cam.Clean()
	}

	w.UpdateChunks()

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
