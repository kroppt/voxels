package world_test

import (
	"testing"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

func TestWorldLoadedChunkCount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc          string
		loadChunks    []chunk.Position
		unloadChunks  []chunk.Position
		expectedCount int
	}{
		{
			desc:          "world starts with no loaded chunks",
			expectedCount: 0,
		},
		{
			desc: "world should load one chunk",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 1,
		},
		{
			desc: "world should load and unload the same chunk",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			unloadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
			},
			expectedCount: 0,
		},
		{
			desc: "world should load two different chunks",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
			},
			expectedCount: 2,
		},
		{
			desc: "world should load 3 chunks and unload two of them",
			loadChunks: []chunk.Position{
				{X: 1, Y: 2, Z: 3},
				{X: 4, Y: 5, Z: 6},
				{X: 7, Y: 8, Z: 9},
			},
			unloadChunks: []chunk.Position{
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
	var actual chunk.Position
	graphicsMod := graphics.FnModule{
		FnLoadChunk: func(ch chunk.Chunk) {
			actual = ch.Position()
		},
	}

	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{})
	worldMod.LoadChunk(chunk.Position{X: 1, Y: 2, Z: 3})
	expected := chunk.Position{X: 1, Y: 2, Z: 3}
	if actual != expected {
		t.Fatalf("expected graphics to receive %v but got %v", expected, actual)
	}
}

func TestWorldUnloadChunkPassesToGraphics(t *testing.T) {
	t.Parallel()
	var actual chunk.Position
	graphicsMod := graphics.FnModule{
		FnUnloadChunk: func(pos chunk.Position) {
			actual = pos
		},
	}
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, &settings.Repository{}, &cache.FnModule{})
	worldMod.LoadChunk(chunk.Position{X: 1, Y: 2, Z: 3})
	worldMod.UnloadChunk(chunk.Position{X: 1, Y: 2, Z: 3})
	expected := chunk.Position{X: 1, Y: 2, Z: 3}
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
		FnGenerateChunk: func(key chunk.Position) chunk.Chunk {
			newChunk := chunk.NewEmpty(key, settingsRepo.GetChunkSize())
			if key == (chunk.Position{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{
					X: 0,
					Y: 0,
					Z: 0,
				}, chunk.BlockTypeDirt)
			}
			return newChunk
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
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})

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
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
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
	worldMod.UnloadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
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
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
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
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, expectedBlockType)
	worldMod.UnloadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
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
	expectSaved := 3
	actualSaved := 0
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			actualSaved++
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, cacheMod)
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeDirt)
	worldMod.LoadChunk(chunk.Position{X: 1, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0}, chunk.BlockTypeDirt)
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 2, Z: 2})
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
	expectSaved := 2
	actualSaved := 0
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	cacheMod := &cache.FnModule{
		FnSave: func(c chunk.Chunk) {
			actualSaved++
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, cacheMod)
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 2, Z: 3})
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeDirt)
	worldMod.UnloadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.Position{X: 2, Y: 4, Z: 0})
	worldMod.LoadChunk(chunk.Position{X: 1, Y: 0, Z: 5})
	worldMod.LoadChunk(chunk.Position{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.Position{X: 8, Y: 5, Z: 5})
	worldMod.SetBlockType(chunk.VoxelCoordinate{X: 8, Y: 5, Z: 5}, chunk.BlockTypeDirt)
	worldMod.Quit()

	if actualSaved != expectSaved {
		t.Fatalf("expected %v chunks to be saved but %v were saved", expectSaved, actualSaved)
	}

}
