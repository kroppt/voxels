package world

import "io/fs"

type Cacheable interface {
	Load(ChunkPos) (*Chunk, bool)
	Save(*Chunk)
}

type Cache struct {
	fs fs.FS
}

func NewCache(fs fs.FS) *Cache {
	return &Cache{fs: fs}
}

func (c *Cache) Load(pos ChunkPos) (*Chunk, bool) {
	_, err := c.fs.Open("chunks")
	if err != nil {
		return nil, false
	}
	return nil, true
}

func (c *Cache) Save(ch *Chunk) {

}
