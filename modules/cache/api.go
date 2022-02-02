package cache

import "github.com/kroppt/voxels/chunk"

type Interface interface {
	Save(chunk.Chunk)
	Load(chunk.ChunkCoordinate) (chunk.Chunk, bool)
	Close()
}

func (m *Module) Save(chunk chunk.Chunk) {
	m.c.save(chunk)
}

func (m *Module) Load(key chunk.ChunkCoordinate) (chunk.Chunk, bool) {
	return m.c.load(key)
}

func (m *Module) Close() {
	m.c.close()
}

type FnModule struct {
	FnSave  func(chunk.Chunk)
	FnLoad  func(chunk.ChunkCoordinate) (chunk.Chunk, bool)
	FnClose func()
}

func (fn *FnModule) Save(chunk chunk.Chunk) {
	if fn.FnSave != nil {
		fn.FnSave(chunk)
	}

}

func (fn *FnModule) Load(pos chunk.ChunkCoordinate) (chunk.Chunk, bool) {
	if fn.FnLoad != nil {
		return fn.FnLoad(pos)
	}
	return chunk.Chunk{}, false
}

func (fn *FnModule) Close() {
	if fn.FnClose != nil {
		fn.FnClose()
	}
}
