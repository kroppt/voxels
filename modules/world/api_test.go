package world_test

import (
	"container/list"
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

func TestWorldLoadedChunkCount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc          string
		loadChunks    []chunk.ChunkCoordinate
		unloadChunks  []chunk.ChunkCoordinate
		expectedCount int
	}{
		{
			desc:          "world starts with no loaded chunks",
			expectedCount: 0,
		},
		{
			desc: "world should load one chunk",
			loadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should load and unload the same chunk",
			loadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
			},
			unloadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should load two different chunks",
			loadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
			},
			expectedCount: 2,
		},
		{
			desc: "world should load 3 chunks and unload two of them",
			loadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
				{X: 7, Y: 8, Z: 9},
			},
			unloadChunks: []chunk.ChunkCoordinate{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
			},
			expectedCount: 1,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			graphicsMod := graphics.FnModule{}
			worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{})
			for _, loadChunk := range tC.loadChunks {
				worldMod.LoadChunk(loadChunk)
			}
			for _, unloadChunk := range tC.unloadChunks {
				worldMod.UnloadChunk(unloadChunk)
			}
			actual := worldMod.CountLoadedChunks()
			if actual != tC.expectedCount {
				t.Fatalf("expected %v chunks to be loaded but got %v", tC.expectedCount, actual)
			}
		})
	}
}

func TestWorldLoadChunkPassesToGraphics(t *testing.T) {
	t.Parallel()
	var actual chunk.ChunkCoordinate
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			actual = ch.Position()
		},
	}

	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3})
	expected := chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3}
	if actual != expected {
		t.Fatalf("expected graphics to receive %v but got %v", expected, actual)
	}
}

func TestWorldUnloadChunkPassesToGraphics(t *testing.T) {
	t.Parallel()
	var actual chunk.ChunkCoordinate
	graphicsMod := graphics.FnModule{
		FnUnloadChunk: func(pos chunk.ChunkCoordinate) {
			actual = pos
		},
	}
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, &settings.FnRepository{}, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3})
	worldMod.UnloadChunk(chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3})
	expected := chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3}
	if actual != expected {
		t.Fatalf("expected graphics to receive %v but got %v", expected, actual)
	}
}

func TestWorldGeneration(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{
					X: 0,
					Y: 0,
					Z: 0,
				}, chunk.BlockTypeDirt)
			}
			return newChunk, list.New()
		},
	}
	expected := chunk.BlockTypeDirt
	var actual chunk.BlockType
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			actual = ch.BlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})

	if actual != expected {
		t.Fatalf("expected to retrieve block type %v but got %v", expected, actual)
	}

}

func TestNewWorldNilSettingsRepo(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	world.New(&graphics.FnModule{}, &world.FnGenerator{}, nil, &cache.FnModule{})
}

func TestNewWorldNilGenerator(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	world.New(&graphics.FnModule{}, nil, &settings.FnRepository{}, &cache.FnModule{})
}

func TestNewWorldNilGraphicsMod(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	world.New(nil, &world.FnGenerator{}, &settings.FnRepository{}, &cache.FnModule{})
}

func TestCannotLoadAlreadyLoadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
}

func TestCannotUnloadNonLoadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.UnloadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
}

func TestCannotSetBlockInUnloadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeAir)
}

func TestCannotGetBlockInUnloadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.GetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
}

func TestValidSetAndGetBlockInWorld(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	expectedBlockType := chunk.BlockTypeDirt
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, expectedBlockType)
	actualBlockType := worldMod.GetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})

	if actualBlockType != expectedBlockType {
		t.Fatalf("expected to receive block type %v but got %v", expectedBlockType, actualBlockType)
	}
}

func TestWorldSavesAfterUnloading(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	expectedBlockType := chunk.BlockTypeDirt
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, cacheMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, expectedBlockType)
	worldMod.UnloadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	actualBlockType := worldMod.GetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})

	if actualBlockType != expectedBlockType {
		t.Fatalf("expected to receive block type %v but got %v", expectedBlockType, actualBlockType)
	}
}

func TestWorldUnloadAllChunksOnQuit(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	expectSaved := 19
	actualSaved := 0
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			actualSaved++
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, cacheMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeDirt)
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 1, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeDirt)
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 2, Z: 2})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 2, Z: 2}, chunk.BlockTypeDirt)

	worldMod.Quit()
	if actualSaved != expectSaved {
		t.Fatalf("expected chunk count to be %v but was %v", expectSaved, actualSaved)
	}
}

func TestWorldClosesCacheOnQuit(t *testing.T) {
	t.Parallel()
	expected := true
	acutal := false
	cacheMod := &cache.FnModule{
		FnClose: func() {
			acutal = true
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, &settings.FnRepository{}, cacheMod)
	worldMod.Quit()
	if acutal != expected {
		t.Fatal("expected quit to call close on cache, but did not")
	}
}

func TestWorldDoesNotSaveChunkIfUnmodified(t *testing.T) {
	t.Parallel()
	expectSaved := 0
	actualSaved := 0
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			actualSaved++
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, cacheMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 2, Z: 3})
	worldMod.Quit()

	if actualSaved != expectSaved {
		t.Fatalf("expected %v chunks to be saved but %v were saved", expectSaved, actualSaved)
	}
}

func TestWorldCallsGraphicsUpdateView(t *testing.T) {
	t.Parallel()
	actual := false
	expected := true
	graphicsMod := graphics.FnModule{
		FnUpdateView: func(map[chunk.ChunkCoordinate]struct{}, mgl.Mat4, chunk.VoxelCoordinate, bool) {
			actual = true
		},
	}
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{},
		Dir: mgl.QuatIdent(),
	})
	if actual != expected {
		t.Fatal("expected world to update view on graphics but didn't")
	}
}

func TestWorldSendsViewMatrix(t *testing.T) {
	t.Parallel()
	// TODO calc this by hand to improve test
	pos := mgl.Vec3{0.5, -1, 2}
	rot := mgl.QuatIdent().Mul(mgl.QuatRotate(mgl.DegToRad(45), mgl.Vec3{1, 1, 1}))
	posNeg := pos.Mul(-1)
	posMat := mgl.Translate3D(posNeg.X(), posNeg.Y(), posNeg.Z())
	expected := mgl.Ident4().Mul4(rot.Inverse().Mat4()).Mul4(posMat)
	var actual mgl.Mat4
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(_ map[chunk.ChunkCoordinate]struct{}, viewMat mgl.Mat4, _ chunk.VoxelCoordinate, _ bool) {
			actual = viewMat
		},
	}
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{})
	worldMod.UpdateView(world.ViewState{
		Pos: pos,
		Dir: rot,
	})
	if actual != expected {
		t.Fatalf("expected graphics to receive view matrix:\n%v but got:\n%v", expected, actual)
	}
}

func TestFrustumCulling(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc       string
		expect     map[chunk.ChunkCoordinate]struct{}
		viewState  world.ViewState
		fov        float64
		far        float64
		near       float64
		width      uint32
		height     uint32
		renderDist uint32
		chunkSize  uint32
	}{
		{
			desc: "simple frustum culling",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: 0, Z: 0}:  {},
				{X: 0, Y: 0, Z: -1}: {},
			},
			viewState: world.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        33.398488467987,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  1,
		},
		{
			desc: "frustum culling large chunks",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: 0, Z: 0}:    {},
				{X: 0, Y: 0, Z: -1}:   {},
				{X: -1, Y: 0, Z: -1}:  {},
				{X: 0, Y: -1, Z: -1}:  {},
				{X: -1, Y: -1, Z: -1}: {},
			},
			viewState: world.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        70,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  3,
		},
		{
			desc: "frustum culling wide angle",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: 0, Z: 0}:    {},
				{X: 0, Y: 0, Z: -1}:   {},
				{X: -1, Y: 0, Z: -1}:  {},
				{X: -1, Y: 1, Z: -1}:  {},
				{X: -1, Y: -1, Z: -1}: {},
				{X: 0, Y: -1, Z: -1}:  {},
				{X: 0, Y: 1, Z: -1}:   {},
				{X: 1, Y: 0, Z: -1}:   {},
				{X: 1, Y: -1, Z: -1}:  {},
				{X: 1, Y: 1, Z: -1}:   {},
			},
			viewState: world.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        89.5,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  1,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actualViewedChunks := map[chunk.ChunkCoordinate]struct{}{}
			graphicsMod := &graphics.FnModule{
				FnUpdateView: func(viewChunks map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4, _ chunk.VoxelCoordinate, _ bool) {
					actualViewedChunks = viewChunks
				},
			}
			settingsMod := settings.FnRepository{
				FnGetFOV: func() float64 {
					return tC.fov
				},
				FnGetFar: func() float64 {
					return tC.far
				},
				FnGetNear: func() float64 {
					return tC.near
				},
				FnGetResolution: func() (uint32, uint32) {
					return tC.width, tC.height
				},
				FnGetRenderDistance: func() uint32 {
					return tC.renderDist
				},
				FnGetChunkSize: func() uint32 {
					return tC.chunkSize
				},
			}
			worldMod := world.New(graphicsMod, &world.FnGenerator{}, settingsMod, &cache.FnModule{})
			worldMod.UpdateView(world.ViewState{
				Pos: tC.viewState.Pos,
				Dir: tC.viewState.Dir,
			})

			if !reflect.DeepEqual(tC.expect, actualViewedChunks) {
				t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", tC.expect, actualViewedChunks)
			}
		})
	}
}

func TestWorldSendsSelectedVoxel(t *testing.T) {
	t.Parallel()
	expected := chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3}
	var actual chunk.VoxelCoordinate
	actualSelected := false
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(_ map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4, selectedVoxel chunk.VoxelCoordinate, selected bool) {
			actual = selectedVoxel
			actualSelected = selected
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
		FnGetRenderDistance: func() uint32 {
			return 0
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3}, chunk.BlockTypeDirt)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.UpdateView(world.ViewState{
		Pos: mgl.Vec3{3.5, 3.5, 4.5},
		Dir: mgl.QuatIdent(),
	})
	if actualSelected != true {
		t.Fatal("expected a voxel to be selected, but one wasn't")
	}
	if actual != expected {
		t.Fatalf("expected to select voxel %v but got %v", expected, actual)
	}
}

func TestWorldSendsSelectedVoxelFromFarAwayChunk(t *testing.T) {
	t.Parallel()
	expected := chunk.VoxelCoordinate{X: 3, Y: 3, Z: -1}
	var actual chunk.VoxelCoordinate
	actualSelected := false
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(_ map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4, selectedVoxel chunk.VoxelCoordinate, selected bool) {
			actual = selectedVoxel
			actualSelected = selected
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
		FnGetRenderDistance: func() uint32 {
			return 1
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 2, Z: 4}, chunk.BlockTypeDirt)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 4, Z: 4}, chunk.BlockTypeDirt)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 4}, chunk.BlockTypeDirt)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 4, Y: 3, Z: 4}, chunk.BlockTypeDirt)
			} else if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: -1}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 3, Z: -1}, chunk.BlockTypeDirt)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 3, Z: -2}, chunk.BlockTypeDirt)
				newChunk.SetAdjacency(chunk.VoxelCoordinate{X: 3, Y: 3, Z: -1}, chunk.AdjacentFront)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 3, Y: 3, Z: -4}, chunk.BlockTypeDirt)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: -1})
	worldMod.UpdateView(world.ViewState{
		Pos: mgl.Vec3{3.5, 3.5, 4.5},
		Dir: mgl.QuatIdent(),
	})
	if actualSelected != true {
		t.Fatal("expected a voxel to be selected, but one wasn't")
	}
	if actual != expected {
		t.Fatalf("expected to select voxel %v but got %v", expected, actual)
	}
}

func TestDeselectAfterMovingAway(t *testing.T) {
	t.Parallel()
	actualSelected := false
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(_ map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4, selectedVoxel chunk.VoxelCoordinate, selected bool) {
			actualSelected = selected
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
		FnGetRenderDistance: func() uint32 {
			return 1
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: -1}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}, chunk.BlockTypeDirt)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	playerMod := player.New(worldMod, settingsRepo)
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: 0.5, Y: 0.5, Z: 0.5})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{Rotation: mgl.QuatIdent()})
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: 0.5, Y: 0.5, Z: 1.5})
	if actualSelected != false {
		t.Fatal("voxel was selected in an unloaded chunk")
	}
}

func TestGraphicsReceivesChunkUpdateOnWorldModification(t *testing.T) {
	t.Parallel()
	var actual chunk.Chunk
	graphicsMod := graphics.FnModule{
		FnUpdateChunk: func(ch chunk.Chunk) {
			actual = ch
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
	}
	expected := chunk.NewChunkEmpty(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}, settingsRepo.GetChunkSize())
	expected.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 2, Z: 3}, chunk.BlockTypeDirt)
	expected.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 3}, chunk.BlockTypeGrass)
	expected.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 3}, chunk.BlockTypeStone)

	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 2, Z: 3}, chunk.BlockTypeDirt)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 3}, chunk.BlockTypeGrass)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 3}, chunk.BlockTypeStone)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected graphics to receive chunk update of %v but got %v", expected, actual)
	}
}

func TestRemoveSelection(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2}, chunk.BlockTypeStone)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphics.FnModule{}, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{2.5, 2.5, 3.5},
		Dir: mgl.QuatIdent(),
	})
	removed := worldMod.RemoveSelection()
	if !removed {
		t.Fatal("expected there to be a selection removal, but there was not")
	}
	expected := chunk.BlockTypeAir
	actual := worldMod.GetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected removed block to be air, but was type # = %v", actual)
	}
}

func TestRemoveSelectionWithoutSelection(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2}, chunk.BlockTypeStone)
			}
			return newChunk, list.New()
		},
	}
	worldMod := world.New(graphics.FnModule{}, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{2.5, 2.5, 1.5},
		Dir: mgl.QuatIdent(),
	})
	removed := worldMod.RemoveSelection()
	if removed {
		t.Fatal("expected there to be no selection removal, but there was")
	}
	expected := chunk.BlockTypeStone
	actual := worldMod.GetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2})

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected non-selected block behind player to be unaffected")
	}
}

func TestSelectAfterRemovingSelection(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 2}, chunk.BlockTypeCorrupted)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 2}, chunk.BlockTypeStone)
				newChunk.SetBlockType(chunk.VoxelCoordinate{X: 2, Y: 2, Z: 3}, chunk.BlockTypeGrass)
			}
			return newChunk, list.New()
		},
	}
	expected := chunk.VoxelCoordinate{
		X: 2,
		Y: 2,
		Z: 2,
	}
	var actual chunk.VoxelCoordinate
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(_ map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4, sv chunk.VoxelCoordinate, _ bool) {
			actual = sv
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{2.5, 2.5, 4.5},
		Dir: mgl.QuatIdent(),
	})
	worldMod.RemoveSelection()
	if actual != expected {
		t.Fatalf("expected new selected voxel to be %v but got %v", expected, actual)
	}
}

func TestSetAdjacenciesAcrossChunksAutomatically(t *testing.T) {
	chunkPos1 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}
	chunkPos2 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: -1}
	var actual1 []float32
	var actual2 []float32
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			if ch.Position() == chunkPos1 {
				actual1 = ch.GetFlatData()
			} else if ch.Position() == chunkPos2 {
				actual2 = ch.GetFlatData()
			}
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	chunk1 := chunk.NewChunkEmpty(chunkPos1, settingsRepo.FnGetChunkSize())
	chunk1.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeCorrupted)
	chunk1.SetAdjacency(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.AdjacentFront)
	expected1 := chunk1.GetFlatData()
	chunk2 := chunk.NewChunkEmpty(chunkPos2, settingsRepo.FnGetChunkSize())
	chunk2.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}, chunk.BlockTypeCorrupted)
	chunk2.SetAdjacency(chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}, chunk.AdjacentBack)
	expected2 := chunk2.GetFlatData()
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			pending := list.New()
			if key == chunkPos1 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeCorrupted))
			} else if key == chunkPos2 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}, chunk.BlockTypeCorrupted))
			}
			return newChunk, pending
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunkPos1)
	worldMod.LoadChunk(chunkPos2)
	if !reflect.DeepEqual(actual1, expected1) {
		t.Fatalf("expected chunk %v to have flat data %v but had %v", chunkPos1, expected1, actual1)
	}
	if !reflect.DeepEqual(actual2, expected2) {
		t.Fatalf("expected chunk %v to have flat data %v but had %v", chunkPos2, expected2, actual2)
	}
}

func TestPendingActionsAlsoUpdatesGraphics(t *testing.T) {
	chunkPos1 := chunk.ChunkCoordinate{X: 1, Y: 0, Z: 0}
	chunkPos2 := chunk.ChunkCoordinate{X: 1, Y: 0, Z: -1}
	var actual []float32
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			if ch.Position() == chunkPos1 {
				actual = ch.GetFlatData()
			}
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	chunk1 := chunk.NewChunkEmpty(chunkPos1, settingsRepo.FnGetChunkSize())
	chunk1.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeCorrupted)
	expected1 := chunk1.GetFlatData()
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			pending := list.New()
			if key == chunkPos1 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeCorrupted))
			} else if key == chunkPos2 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: -1}, chunk.BlockTypeCorrupted))
			}
			return newChunk, pending
		},
	}
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{1.5, 0.5, -2.5},
		Dir: mgl.Quat{
			W: 0,
			V: [3]float64{0, -1, 0},
		},
	})
	worldMod.LoadChunk(chunkPos1)
	worldMod.LoadChunk(chunkPos2)
	worldMod.RemoveSelection()

	if !reflect.DeepEqual(actual, expected1) {
		t.Fatalf("expected chunk %v to have flat data %v but had %v", chunkPos1, expected1, actual)
	}
}

func TestUpdateViewAfterUnload(t *testing.T) {
	chunkPos1 := chunk.ChunkCoordinate{X: 1, Y: 0, Z: 0}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			pending := list.New()
			if key == chunkPos1 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeCorrupted))
			}
			return newChunk, pending
		},
	}
	worldMod := world.New(graphics.FnModule{}, testGen, settingsRepo, &cache.FnModule{})
	worldMod.LoadChunk(chunkPos1)
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{1.5, 0.5, -2.5},
		Dir: mgl.Quat{
			W: 0,
			V: [3]float64{0, -1, 0},
		},
	})
	worldMod.UnloadChunk(chunkPos1)
	worldMod.RemoveSelection()
}

func TestGraphicsExpectedLoadsAndUpdates(t *testing.T) {
	expectedLoads := 1
	expectedUpdates := 1
	actualLoads := 0
	actualUpdates := 0
	pos := chunk.ChunkCoordinate{X: 1, Y: 0, Z: 0}
	pos2 := chunk.ChunkCoordinate{X: 1, Y: 0, Z: -1}
	graphicMod := graphics.FnModule{
		FnLoadChunk: func(c chunk.Chunk) {
			if c.Position() == pos {
				actualLoads++
			}
		},
		FnUpdateChunk: func(c chunk.Chunk) {
			if c.Position() == pos {
				actualUpdates++
			}
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			newChunk := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
			pending := list.New()
			if key == pos {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeCorrupted))
			} else if key == pos2 {
				pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: -1}, chunk.BlockTypeCorrupted))
			}
			return newChunk, pending
		},
	}
	worldMod := world.New(graphicMod, testGen, settingsRepo, &cache.FnModule{})
	worldMod.UpdateView(world.ViewState{
		Pos: [3]float64{1.5, 0.5, -2.5},
		Dir: mgl.Quat{
			W: 0,
			V: [3]float64{0, -1, 0},
		},
	})
	worldMod.LoadChunk(pos)
	worldMod.LoadChunk(pos2)

	if expectedLoads != actualLoads {
		t.Fatalf("expected %v loads for %v but got %v loads", expectedLoads, pos, actualLoads)
	}
	if expectedUpdates != actualUpdates {
		t.Fatalf("expected %v updates for %v but got %v updates", expectedUpdates, pos, actualUpdates)
	}
}

// TODO test that world saves pending actions on quit by loading each chunk with pending actions,
// applying the actions, and then saving those chunks
