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
	}
	cacheMod := cache.New(afero.NewMemMapFs(), settingsRepo)
	chPos := chunk.Position{X: 1, Y: 2, Z: 3}
	vPos := chunk.VoxelCoordinate{X: 1, Y: 2, Z: 3}
	testChunk := chunk.NewEmpty(chPos, settingsRepo.GetChunkSize())
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
	chPos := chunk.Position{X: -1, Y: -1, Z: -1}
	testChunk := chunk.NewEmpty(chPos, settingsRepo.GetChunkSize())
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
	chPos1 := chunk.Position{X: -1, Y: -1, Z: -1}
	testChunk1 := chunk.NewEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.LightFront, 5)
	chPos2 := chunk.Position{X: 0, Y: 0, Z: 0}
	testChunk2 := chunk.NewEmpty(chPos2, settingsRepo.GetChunkSize())
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
	chPos1 := chunk.Position{X: -1, Y: -1, Z: -1}
	testChunk1 := chunk.NewEmpty(chPos1, settingsRepo.GetChunkSize())
	testChunk1.SetAdjacency(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.AdjacentAll)
	testChunk1.SetLighting(chunk.VoxelCoordinate{X: -1, Y: -1, Z: -1}, chunk.LightFront, 5)
	chPos2 := chunk.Position{X: 0, Y: 0, Z: 0}
	testChunk2 := chunk.NewEmpty(chPos2, settingsRepo.GetChunkSize())
	testChunk2.SetAdjacency(chunk.VoxelCoordinate{X: 2, Y: 3, Z: 4}, chunk.AdjacentBottom)
	testChunk2.SetLighting(chunk.VoxelCoordinate{X: 5, Y: 6, Z: 7}, chunk.LightLeft, 3)
	chPos3 := chunk.Position{X: 5, Y: 5, Z: 5}
	testChunk3 := chunk.NewEmpty(chPos3, settingsRepo.GetChunkSize())
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
