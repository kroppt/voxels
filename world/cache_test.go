package world_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kroppt/voxels/world"
)

func compareFlatData(d1, d2 []float32) bool {
	if len(d1) != len(d2) {
		return false
	}
	for i := 0; i < len(d1); i++ {
		if d1[i] != d2[i] {
			return false
		}
	}
	return true
}

func TestCacheInMemory(t *testing.T) {
	cache, err := world.NewCache("test_meta", "test_data", 2)
	if err != nil {
		panic(fmt.Sprintf("failed to init cache: %v", err))
	}
	ch := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	}, world.FlatWorldGenerator{})
	cache.Save(ch)
	ch2, loaded := cache.Load(world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	})
	if !loaded {
		t.Fatal("expected chunk to be loaded but was not")
	}
	t.Logf("ch: %v", ch)
	t.Logf("ch2: %v", ch2)
	cache.Destroy()
	if !compareFlatData(ch.GetFlatData(), ch2.GetFlatData()) {
		t.Fatalf("loaded data not same as saved data")
	}
	err = os.Remove("test_data")
	if err != nil {
		t.Fatalf("failed to remove test_data")
	}
	err = os.Remove("test_meta")
	if err != nil {
		t.Fatalf("failed to remove test_meta")
	}
}

func TestCacheInFile(t *testing.T) {
	cache, err := world.NewCache("test_meta", "test_data", 1)
	if err != nil {
		panic(fmt.Sprintf("failed to init cache: %v", err))
	}
	ch := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	}, world.FlatWorldGenerator{})
	cache.Save(ch)
	ch2, loaded := cache.Load(world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	})
	if !loaded {
		t.Fatal("expected chunk to be loaded but was not")
	}
	t.Logf("ch: %v", ch)
	t.Logf("ch2: %v", ch2)
	cache.Destroy()
	if !compareFlatData(ch.GetFlatData(), ch2.GetFlatData()) {
		t.Fatalf("loaded data not same as saved data")
	}
	err = os.Remove("test_data")
	if err != nil {
		t.Fatalf("failed to remove test_data")
	}
	err = os.Remove("test_meta")
	if err != nil {
		t.Fatalf("failed to remove test_meta")
	}
}

func TestCacheGetNumChunksMeta(t *testing.T) {
	cache, err := world.NewCache("test_meta", "test_data", 2)
	if err != nil {
		panic(fmt.Sprintf("failed to init cache: %v", err))
	}
	ch := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	}, world.FlatWorldGenerator{})
	cache.Save(ch)
	ch2 := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 1,
		Y: 0,
		Z: 2,
	}, world.FlatWorldGenerator{})
	cache.Save(ch2)
	chunksInFile, success := cache.GetNumChunksInFile()
	if !success {
		t.Fatal("failed to check meta file for # chunks")
	}
	if chunksInFile != 2 {
		t.Fatalf("expected 2 chunks in file but got %v", chunksInFile)
	}

	cache.Destroy()
	err = os.Remove("test_data")
	if err != nil {
		t.Fatalf("failed to remove test_data")
	}
	err = os.Remove("test_meta")
	if err != nil {
		t.Fatalf("failed to remove test_meta")
	}
}

func TestCache2Chunks(t *testing.T) {
	cache, err := world.NewCache("test_meta", "test_data", 2)
	if err != nil {
		panic(fmt.Sprintf("failed to init cache: %v", err))
	}
	ch := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	}, world.FlatWorldGenerator{})
	cache.Save(ch)
	ch2 := world.NewChunk(world.ChunkSize, world.ChunkPos{
		X: 1,
		Y: 0,
		Z: 2,
	}, world.FlatWorldGenerator{})
	cache.Save(ch2)
	chLoaded, loaded := cache.Load(world.ChunkPos{
		X: 0,
		Y: 0,
		Z: 0,
	})
	ch2Loaded, loaded2 := cache.Load(world.ChunkPos{
		X: 1,
		Y: 0,
		Z: 2,
	})
	if !loaded {
		t.Fatal("expected chunk1 to be loaded but was not")
	}
	if !loaded2 {
		t.Fatal("expected chunk2 to be loaded but was not")
	}
	cache.Destroy()
	if !compareFlatData(ch.GetFlatData(), chLoaded.GetFlatData()) {
		t.Fatalf("loaded data not same as saved data")
	}
	if !compareFlatData(ch2.GetFlatData(), ch2Loaded.GetFlatData()) {
		t.Fatalf("loaded data not same as saved data")
	}
	err = os.Remove("test_data")
	if err != nil {
		t.Fatalf("failed to remove test_data")
	}
	err = os.Remove("test_meta")
	if err != nil {
		t.Fatalf("failed to remove test_meta")
	}
}

func TestCacheManyChunks(t *testing.T) {
	cache, err := world.NewCache("test_meta", "test_data", 5)
	if err != nil {
		panic(fmt.Sprintf("failed to init cache: %v", err))
	}
	nChunks := 15
	for i := 0; i < nChunks; i++ {
		ch := world.NewChunk(world.ChunkSize, world.ChunkPos{
			X: i,
			Y: i,
			Z: i,
		}, world.FlatWorldGenerator{})
		cache.Save(ch)
	}
	resultingNumChunks, success := cache.GetNumChunksInFile()
	if !success || resultingNumChunks != int32(nChunks) {
		t.Fatalf("expected %v chunks but got %v", nChunks, resultingNumChunks)
	}
	cache.Destroy()
	err = os.Remove("test_data")
	if err != nil {
		t.Fatalf("failed to remove test_data")
	}
	err = os.Remove("test_meta")
	if err != nil {
		t.Fatalf("failed to remove test_meta")
	}
}
