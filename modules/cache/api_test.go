package cache_test

import (
	"reflect"
	"testing"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

func TestNewCacheNotNil(t *testing.T) {
	t.Parallel()
	cacheMod := cache.New(afero.NewMemMapFs(), settings.FnRepository{})
	if cacheMod == nil {
		t.Fatal("expected new cache to not be nil, but it was")
	}
}

func TestCacheSettingsRepoNonNil(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	cacheMod := cache.New(afero.NewMemMapFs(), nil)
	if cacheMod == nil {
		t.Fatal("expected new cache to not be nil, but it was")
	}
}

func TestCacheReadAndWriteSimpleChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 1
		},
		FnGetRegionSize: func() uint32 {
			return 1
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos := chunk.ChunkCoordinate{X: 1, Y: 2, Z: 3}
	vPos := chunk.VoxelCoordinate{X: 1, Y: 2, Z: 3}
	testChunk := chunk.NewChunkEmpty(chPos, settingsRepo.GetChunkSize())
	testChunk.SetAdjacency(vPos, chunk.AdjacentAll)
	testChunk.SetLighting(vPos, chunk.LightFront, 5)
	expectedData := testChunk.GetFlatData()
	cacheMod.Save(testChunk)
	loadedChunk, loaded := cacheMod.Load(chPos)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData := loadedChunk.GetFlatData()
	if !reflect.DeepEqual(actualData, expectedData) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData, chPos, actualData)
	}
}

func TestCacheReadAndWriteComplexChunk(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 10
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos := chunk.ChunkCoordinate{X: -1, Y: -1, Z: -1}
	testChunk := chunk.NewChunkEmpty(chPos, settingsRepo.GetChunkSize())
	testChunk.SetAdjacency(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.AdjacentAll)
	testChunk.SetLighting(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.LightFront, 5)
	testChunk.SetBlockType(chunk.VoxelCoordinate{X: -2, Y: -2, Z: -2}, chunk.BlockTypeAir)
	testChunk.SetBlockType(chunk.VoxelCoordinate{X: -1, Y: -2, Z: -2}, chunk.BlockTypeDirt)
	testChunk.SetBlockType(chunk.VoxelCoordinate{X: -2, Y: -2, Z: -1}, chunk.BlockTypeDirt)
	testChunk.SetBlockType(chunk.VoxelCoordinate{X: -10, Y: -10, Z: -10}, chunk.BlockTypeDirt)

	expectedData := testChunk.GetFlatData()
	cacheMod.Save(testChunk)
	loadedChunk, loaded := cacheMod.Load(chPos)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData := loadedChunk.GetFlatData()
	if !reflect.DeepEqual(actualData, expectedData) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData, chPos, actualData)
	}
}

func TestCacheReadAndWriteTwoChunks(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 10
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos1 := chunk.ChunkCoordinate{X: -1, Y: -1, Z: -1}
	testChunk1 := chunk.NewChunkEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.LightFront, 5)
	chPos2 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}
	testChunk2 := chunk.NewChunkEmpty(chPos2, settingsRepo.GetChunkSize())
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 4}, chunk.AdjacentBottom)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 5, Y: 6, Z: 7}, chunk.LightLeft, 3)

	expectedData1 := testChunk1.GetFlatData()
	expectedData2 := testChunk2.GetFlatData()
	cacheMod.Save(testChunk1)
	cacheMod.Save(testChunk2)
	loadedChunk1, loaded := cacheMod.Load(chPos1)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData1 := loadedChunk1.GetFlatData()
	loadedChunk2, loaded := cacheMod.Load(chPos2)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData2 := loadedChunk2.GetFlatData()
	if !reflect.DeepEqual(actualData1, expectedData1) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData1, chPos1, actualData1)
	}
	if !reflect.DeepEqual(actualData2, expectedData2) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData2, chPos2, actualData2)
	}
}

func TestCacheReadWriteAndOverwriteThreeChunks(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 10
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos1 := chunk.ChunkCoordinate{X: -1, Y: -1, Z: -1}
	testChunk1 := chunk.NewChunkEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.LightFront, 5)
	chPos2 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}
	testChunk2 := chunk.NewChunkEmpty(chPos2, settingsRepo.GetChunkSize())
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 4}, chunk.AdjacentBottom)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 5, Y: 6, Z: 7}, chunk.LightLeft, 3)
	chPos3 := chunk.ChunkCoordinate{X: 5, Y: 5, Z: 5}
	testChunk3 := chunk.NewChunkEmpty(chPos3, settingsRepo.GetChunkSize())
	testChunk3.SetAdjacency(chunk.VoxelCoordinate{X: 51, Y: 51, Z: 51}, chunk.AdjacentTop)
	testChunk3.SetLighting(chunk.VoxelCoordinate{X: 55, Y: 55, Z: 52}, chunk.LightRight, 1)

	cacheMod.Save(testChunk1)
	cacheMod.Save(testChunk2)
	cacheMod.Save(testChunk3)
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 4}, chunk.AdjacentX)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 7, Y: 7, Z: 7}, chunk.LightBottom, 0)
	cacheMod.Save(testChunk2)
	expectedData1 := testChunk1.GetFlatData()
	expectedData2 := testChunk2.GetFlatData()
	expectedData3 := testChunk3.GetFlatData()
	loadedChunk1, loaded := cacheMod.Load(chPos1)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData1 := loadedChunk1.GetFlatData()
	loadedChunk2, loaded := cacheMod.Load(chPos2)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData2 := loadedChunk2.GetFlatData()
	loadedChunk3, loaded := cacheMod.Load(chPos3)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData3 := loadedChunk3.GetFlatData()
	if !reflect.DeepEqual(actualData1, expectedData1) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData1, chPos1, actualData1)
	}
	if !reflect.DeepEqual(actualData2, expectedData2) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData2, chPos2, actualData2)
	}
	if !reflect.DeepEqual(actualData3, expectedData3) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData3, chPos3, actualData3)
	}
}

func TestCacheReadAndWriteChunksInDifferentRegions(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 2
		},
		FnGetRegionSize: func() uint32 {
			return 5
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos1 := chunk.ChunkCoordinate{X: -100, Y: -100, Z: 345}
	testChunk1 := chunk.NewChunkEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: -200, Y: -200, Z: 690}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: -200, Y: -200, Z: 690}, chunk.LightFront, 5)
	chPos2 := chunk.ChunkCoordinate{X: 66, Y: -70, Z: 0}
	testChunk2 := chunk.NewChunkEmpty(chPos2, settingsRepo.GetChunkSize())
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 132, Y: -140, Z: 0}, chunk.AdjacentBottom)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 133, Y: -139, Z: 1}, chunk.LightLeft, 3)

	expectedData1 := testChunk1.GetFlatData()
	expectedData2 := testChunk2.GetFlatData()
	cacheMod.Save(testChunk1)
	cacheMod.Save(testChunk2)
	loadedChunk1, loaded := cacheMod.Load(chPos1)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData1 := loadedChunk1.GetFlatData()
	loadedChunk2, loaded := cacheMod.Load(chPos2)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData2 := loadedChunk2.GetFlatData()
	if !reflect.DeepEqual(actualData1, expectedData1) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData1, chPos1, actualData1)
	}
	if !reflect.DeepEqual(actualData2, expectedData2) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData2, chPos2, actualData2)
	}
}

func TestCacheReadAndWriteChunksInSameRegion(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
		FnGetRegionSize: func() uint32 {
			return 5
		},
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos1 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}
	testChunk1 := chunk.NewChunkEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: 3, Y: 3, Z: 3}, chunk.LightFront, 5)
	chPos2 := chunk.ChunkCoordinate{X: 0, Y: 0, Z: 1}
	testChunk2 := chunk.NewChunkEmpty(chPos2, settingsRepo.GetChunkSize())
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 3, Y: 3, Z: 5}, chunk.AdjacentBottom)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 3, Y: 3, Z: 5}, chunk.LightLeft, 3)

	expectedData1 := testChunk1.GetFlatData()
	expectedData2 := testChunk2.GetFlatData()
	cacheMod.Save(testChunk1)
	cacheMod.Save(testChunk2)
	loadedChunk1, loaded := cacheMod.Load(chPos1)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData1 := loadedChunk1.GetFlatData()
	loadedChunk2, loaded := cacheMod.Load(chPos2)
	if !loaded {
		t.Fatal("failed to load")
	}
	actualData2 := loadedChunk2.GetFlatData()
	if !reflect.DeepEqual(actualData1, expectedData1) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData1, chPos1, actualData1)
	}
	if !reflect.DeepEqual(actualData2, expectedData2) {
		t.Fatalf("expected to retrieve data %v for chunk %v but instead got %v", expectedData2, chPos2, actualData2)
	}
}

func TestCacheReadAndWriteAllInRenderDistance(t *testing.T) {
	t.Parallel()
	settingsRepo := settings.FnRepository{
		FnGetChunkSize: func() uint32 {
			return 5
		},
		FnGetRegionSize: func() uint32 {
			return 5
		},
	}

	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	loadedChunks := map[chunk.ChunkCoordinate]*chunk.Chunk{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				key := chunk.ChunkCoordinate{X: int32(x), Y: int32(y), Z: int32(z)}
				c := chunk.NewChunkEmpty(key, settingsRepo.GetChunkSize())
				v := chunk.VoxelCoordinate{
					X: int32(x * int(settingsRepo.GetChunkSize())),
					Y: int32(y * int(settingsRepo.GetChunkSize())),
					Z: int32(z * int(settingsRepo.GetChunkSize())),
				}
				c.SetBlockType(v, chunk.BlockTypeDirt)
				c.SetAdjacency(v, chunk.AdjacentBottom)
				c.SetLighting(v, chunk.LightTop, 3)
				loadedChunks[key] = &c
				cacheMod.Save(c)
			}
		}
	}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				key := chunk.ChunkCoordinate{X: int32(x), Y: int32(y), Z: int32(z)}
				c, loaded := cacheMod.Load(key)
				if !loaded {
					t.Fatal("failed to load")
				}
				expectedData := loadedChunks[key].GetFlatData()
				actualData := c.GetFlatData()
				if !reflect.DeepEqual(actualData, expectedData) {
					t.Fatalf("expected chunk at %v to have data: %v, but had data: %v", key, expectedData, actualData)
				}
			}
		}
	}
}
