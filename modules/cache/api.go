package cache

import "github.com/kroppt/voxels/chunk"

type Interface interface {
	Save(chunk.Chunk)
	Load(chunk.Position) (chunk.Chunk, bool)
}

func (m *Module) Save(chunk chunk.Chunk) {
	m.c.save(chunk)
}

func (m *Module) Load(key chunk.Position) (chunk.Chunk, bool) {
	return m.c.load(key)
}

type FnModule struct {
	FnSave func(chunk.Chunk)
	FnLoad func(chunk.Position) (chunk.Chunk, bool)
}

func (fn *FnModule) Save(chunk chunk.Chunk) {
	if fn.FnSave != nil {
		fn.FnSave(chunk)
	}

}

func (fn *FnModule) Load(pos chunk.Position) (chunk.Chunk, bool) {
	if fn.FnLoad != nil {
		return fn.FnLoad(pos)
	}
	return chunk.Chunk{}, false
}
