package world_test

import (
	"io/fs"
	"testing"

	"github.com/kroppt/voxels/world"
)

type MemMapFS struct {
	files map[string]MemFile
}

func NewMemMapFS() *MemMapFS {
	var fs MemMapFS
	fs.files = make(map[string]MemFile)
	return &fs
}

func (m MemMapFS) Open(name string) (fs.File, error) {
	if file, ok := m.files[name]; ok {
		return file, nil
	}
	return nil, fs.ErrNotExist
}

type MemFile struct {
}

func (m MemFile) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (m MemFile) Read([]byte) (int, error) {
	return 0, nil
}

func (m MemFile) Close() error {
	return nil
}

func (m MemFile) Write([]byte) (int, error) {
	return 0, nil
}

func TestNewCache(t *testing.T) {
	cache := world.NewCache(MemMapFS{})
	if cache == nil {
		t.Fatal("expected valid cache")
	}
}

func TestCacheLoad(t *testing.T) {
	t.Run("empty fs load", func(t *testing.T) {
		t.Parallel()
		cache := world.NewCache(MemMapFS{})
		_, loaded := cache.Load(world.ChunkPos{0, 0, 0})
		if loaded {
			t.Fatal("expected load to fail")
		}
	})
	t.Run("has entries, but not the correct one to load", func(t *testing.T) {
		t.Parallel()
		fs := MemMapFS{
			files: map[string]MemFile{
				"1": {},
			},
		}
		cache := world.NewCache(fs)
		_, loaded := cache.Load(world.ChunkPos{0, 0, 0})
		if loaded {
			t.Fatal("expected load to fail")
		}
	})
	t.Run("has entries and loads correct one", func(t *testing.T) {
		t.Parallel()
		fs := MemMapFS{
			files: map[string]MemFile{
				"chunks": {},
			},
		}
		cache := world.NewCache(fs)
		_, loaded := cache.Load(world.ChunkPos{0, 0, 0})
		if !loaded {
			t.Fatal("expected load to work")
		}
	})
}
