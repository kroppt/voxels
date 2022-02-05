package world_test

import (
	"container/list"
	"reflect"
	"testing"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
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
			worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
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

	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
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
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, &settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
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
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})

	if actual != expected {
		t.Fatalf("expected to retrieve block type %v but got %v", expected, actual)
	}

}

func TestNewWorldRequiredModules(t *testing.T) {
	t.Parallel()
	t.Run("nil settings repo", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		world.New(&graphics.FnModule{}, &world.FnGenerator{}, nil, &cache.FnModule{}, &view.FnModule{})
	})
	t.Run("nil generator", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		world.New(&graphics.FnModule{}, nil, &settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
	})
	t.Run("nil graphics mod", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		world.New(nil, &world.FnGenerator{}, &settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
	})
	t.Run("nil view mod", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		world.New(&graphics.FnModule{}, &world.FnGenerator{}, &settings.FnRepository{}, &cache.FnModule{}, nil)
	})
}

func TestCannotLoadAlreadyLoadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{}, &view.FnModule{})
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
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.UnloadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
}

func TestCannotGetBlockInUnloadedChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod.GetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
}

func TestWorldSavesAfterUnloading(t *testing.T) {
	t.Parallel()
	saved := false
	loaded := false
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			saved = true
		},
		FnLoad: func(cc chunk.ChunkCoordinate) (chunk.Chunk, bool) {
			loaded = true
			return chunk.Chunk{}, false
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 2
		},
	}
	gen := world.FnGenerator{
		FnGenerateChunk: func(chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{}, settingsRepo.GetChunkSize())
			l := list.New()
			l.PushBackList(c.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}, chunk.BlockTypeStone))
			return c, l
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &gen, settingsRepo, cacheMod, &view.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.RemoveBlock(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1})
	worldMod.UnloadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})

	if !saved {
		t.Fatal("expected chunk to be saved, but it wasn't")
	}
	if !loaded {
		t.Fatal("expected chunk to be loaded, but it wasn't")
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
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, &settings.FnRepository{}, cacheMod, &view.FnModule{})
	worldMod.Quit()
	if acutal != expected {
		t.Fatal("expected quit to call close on cache, but did not")
	}
}

func TestWorldDoesNotSaveChunkIfUnmodified(t *testing.T) {
	t.Parallel()
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			t.Fatal("chunk was saved when none were expected to")
		},
	}
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{}, settings.FnRepository{}, cacheMod, &view.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 2, Z: 3})
	worldMod.Quit()
}

func TestUpdateGraphicsOnRemove(t *testing.T) {
	t.Parallel()
	expected := chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3}
	var actual chunk.ChunkCoordinate
	graphicsMod := graphics.FnModule{
		FnUpdateChunk: func(ch chunk.Chunk) {
			actual = ch.Position()
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
	}
	worldMod := world.New(graphicsMod, &world.FnGenerator{}, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3})
	worldMod.RemoveBlock(chunk.VoxelCoordinate{X: 1, Y: 2, Z: 3})
	if actual != expected {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}

func TestRemoveBlockPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	worldMod := world.New(graphics.FnModule{}, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})
	worldMod.RemoveBlock(chunk.VoxelCoordinate{})
}

func TestLoadChunkUpdateSelection(t *testing.T) {
	t.Parallel()
	updated := false
	viewMod := view.FnModule{
		FnUpdateSelection: func() {
			updated = true
		},
	}
	worldMod := world.New(graphics.FnModule{}, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{}, &viewMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{})
	if !updated {
		t.Fatal("expected selection to be updated, but it was not")
	}
}

func TestUnloadChunkUpdateSelection(t *testing.T) {
	t.Parallel()
	updated := false
	viewMod := view.FnModule{
		FnUpdateSelection: func() {
			updated = true
		},
	}
	worldMod := world.New(graphics.FnModule{}, &world.FnGenerator{}, settings.FnRepository{}, &cache.FnModule{}, &viewMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{})
	worldMod.UnloadChunk(chunk.ChunkCoordinate{})
	if !updated {
		t.Fatal("expected selection to be updated, but it was not")
	}
}

func TestLoadChunkAddNode(t *testing.T) {
	t.Parallel()
	expected := map[chunk.VoxelCoordinate]struct{}{}
	expected[chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}] = struct{}{}
	expected[chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}] = struct{}{}

	actual := map[chunk.VoxelCoordinate]struct{}{}
	viewMod := view.FnModule{
		FnAddNode: func(vc chunk.VoxelCoordinate) {
			actual[vc] = struct{}{}
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 2
		},
	}
	gen := world.FnGenerator{
		FnGenerateChunk: func(chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{}, settingsRepo.GetChunkSize())
			l := list.New()
			l.PushBackList(c.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeStone))
			l.PushBackList(c.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}, chunk.BlockTypeStone))
			return c, l
		},
	}
	worldMod := world.New(graphics.FnModule{}, &gen, settingsRepo, &cache.FnModule{}, &viewMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{})
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected add nodes %v but got %v", expected, actual)
	}
}

func TestUnloadChunkRemoveNode(t *testing.T) {
	t.Parallel()
	expected := map[chunk.VoxelCoordinate]struct{}{}
	expected[chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}] = struct{}{}
	expected[chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}] = struct{}{}

	actual := map[chunk.VoxelCoordinate]struct{}{}
	viewMod := view.FnModule{
		FnRemoveNode: func(vc chunk.VoxelCoordinate) {
			actual[vc] = struct{}{}
		},
	}
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 2
		},
	}
	gen := world.FnGenerator{
		FnGenerateChunk: func(chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{}, settingsRepo.GetChunkSize())
			l := list.New()
			l.PushBackList(c.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeStone))
			l.PushBackList(c.SetBlockType(chunk.VoxelCoordinate{X: 1, Y: 1, Z: 1}, chunk.BlockTypeStone))
			return c, l
		},
	}
	worldMod := world.New(graphics.FnModule{}, &gen, settingsRepo, &cache.FnModule{}, &viewMod)
	worldMod.LoadChunk(chunk.ChunkCoordinate{})
	worldMod.UnloadChunk(chunk.ChunkCoordinate{})
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected remove nodes %v but got %v", expected, actual)
	}
}

func TestSetAdjacenciesAcrossChunksAutomatically(t *testing.T) {
	t.Parallel()
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
	chunk1 := chunk.NewChunkEmpty(chunkPos1, settingsRepo.GetChunkSize())
	chunk1.SetBlockType(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.BlockTypeCorrupted)
	chunk1.SetAdjacency(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0}, chunk.AdjacentFront)
	expected1 := chunk1.GetFlatData()
	chunk2 := chunk.NewChunkEmpty(chunkPos2, settingsRepo.GetChunkSize())
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
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{}, &view.FnModule{})
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
	t.Parallel()
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
	chunk1 := chunk.NewChunkEmpty(chunkPos1, settingsRepo.GetChunkSize())
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
	worldMod := world.New(graphicsMod, testGen, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	worldMod.LoadChunk(chunkPos1)
	worldMod.LoadChunk(chunkPos2)
	worldMod.RemoveBlock(chunk.VoxelCoordinate{X: 1, Y: 0, Z: -1})

	if !reflect.DeepEqual(actual, expected1) {
		t.Fatalf("expected chunk %v to have flat data %v but had %v", chunkPos1, expected1, actual)
	}
}

func TestGraphicsExpectedLoadsAndUpdates(t *testing.T) {
	t.Parallel()
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
	worldMod := world.New(graphicMod, testGen, settingsRepo, &cache.FnModule{}, &view.FnModule{})
	worldMod.LoadChunk(pos)
	worldMod.LoadChunk(pos2)

	if expectedLoads != actualLoads {
		t.Fatalf("expected %v loads for %v but got %v loads", expectedLoads, pos, actualLoads)
	}
	if expectedUpdates != actualUpdates {
		t.Fatalf("expected %v updates for %v but got %v updates", expectedUpdates, pos, actualUpdates)
	}
}

func TestWorldUnloadAllChunksOnQuit(t *testing.T) {
	t.Parallel()
	expectSaved := 18
	actualSaved := 0
	cacheMod := &cache.FnModule{
		FnSave: func(chunk.Chunk) {
			actualSaved++
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
			pending.PushBackList(newChunk.SetBlockType(chunk.VoxelCoordinate{X: key.X, Y: key.Y, Z: key.Z}, chunk.BlockTypeCorrupted))
			return newChunk, pending
		},
	}
	worldMod := world.New(&graphics.FnModule{}, testGen, settingsRepo, cacheMod, &view.FnModule{})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.RemoveBlock(chunk.VoxelCoordinate{X: 0, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 1, Y: 0, Z: 0})
	worldMod.RemoveBlock(chunk.VoxelCoordinate{X: 1, Y: 0, Z: 0})
	worldMod.LoadChunk(chunk.ChunkCoordinate{X: 0, Y: 2, Z: 2})
	worldMod.Quit()

	if actualSaved != expectSaved {
		t.Fatalf("expected chunk count to be %v but was %v", expectSaved, actualSaved)
	}
}

func BenchmarkWorldLoadUnload(b *testing.B) {
	chPos := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}
	ch := chunk.NewChunkEmpty(chPos, 5)
	ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
		ch.SetBlockType(vc, chunk.BlockTypeDirt)
	})
	actions := list.New()
	worldMod := world.New(&graphics.FnModule{}, &world.FnGenerator{
		FnGenerateChunk: func(coord chunk.ChunkCoordinate) (chunk.Chunk, *list.List) {
			return ch, actions
		},
	}, settings.FnRepository{}, &cache.FnModule{}, &view.FnModule{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		worldMod.LoadChunk(chPos)
		worldMod.UnloadChunk(chPos)
	}
}
