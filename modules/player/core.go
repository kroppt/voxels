package player

import (
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod  graphics.Interface
	settingsMod  settings.Interface
	chunkSize    uint32
	lastChunkPos chunkPos
}

func (c *core) init() {
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	for x := int32(-renderDistance); x <= renderDistance; x++ {
		for y := int32(-renderDistance); y <= renderDistance; y++ {
			for z := int32(-renderDistance); z <= renderDistance; z++ {
				c.graphicsMod.ShowChunk(graphics.ChunkEvent{PositionX: x, PositionY: y, PositionZ: z})
			}
		}
	}
}

func (c *core) playerToChunkPosition(pos PositionEvent) chunkPos {
	x, y, z := pos.X, pos.Y, pos.Z
	chunkSize := int32(c.chunkSize)
	if pos.X < 0 {
		x++
	}
	if pos.Y < 0 {
		y++
	}
	if pos.Z < 0 {
		z++
	}
	x /= chunkSize
	y /= chunkSize
	z /= chunkSize
	if pos.X < 0 {
		x--
	}
	if pos.Y < 0 {
		y--
	}
	if pos.Z < 0 {
		z--
	}
	return chunkPos{x, y, z}
}

type chunkPos struct {
	x int32
	y int32
	z int32
}

// chunkRange is the range of chunks between Min and Max.
type chunkRange struct {
	Min chunkPos
	Max chunkPos
}

// forEach executes the given function on every position in the this ChunkRange.
// The return of fn indices whether to stop iterating
func (rng chunkRange) forEach(fn func(pos chunkPos) bool) {
	for x := rng.Min.x; x <= rng.Max.x; x++ {
		for y := rng.Min.y; y <= rng.Max.y; y++ {
			for z := rng.Min.z; z <= rng.Max.z; z++ {
				stop := fn(chunkPos{x: x, y: y, z: z})
				if stop {
					return
				}
			}
		}
	}
}

// contains returns whether this ChunkRange contains the given pos.
func (rng chunkRange) contains(pos chunkPos) bool {
	if pos.x < rng.Min.x || pos.x > rng.Max.x {
		return false
	}
	if pos.y < rng.Min.y || pos.y > rng.Max.y {
		return false
	}
	if pos.z < rng.Min.z || pos.z > rng.Max.z {
		return false
	}
	return true
}

func (c *core) updatePosition(posEvent PositionEvent) {
	newChunkPos := c.playerToChunkPosition(posEvent)
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	old := chunkRange{
		Min: chunkPos{
			x: c.lastChunkPos.x - renderDistance,
			y: c.lastChunkPos.y - renderDistance,
			z: c.lastChunkPos.z - renderDistance,
		},
		Max: chunkPos{
			x: c.lastChunkPos.x + renderDistance,
			y: c.lastChunkPos.y + renderDistance,
			z: c.lastChunkPos.z + renderDistance,
		},
	}
	new := chunkRange{
		Min: chunkPos{
			x: newChunkPos.x - renderDistance,
			y: newChunkPos.y - renderDistance,
			z: newChunkPos.z - renderDistance,
		},
		Max: chunkPos{
			x: newChunkPos.x + renderDistance,
			y: newChunkPos.y + renderDistance,
			z: newChunkPos.z + renderDistance,
		},
	}
	new.forEach(func(pos chunkPos) bool {
		if !old.contains(pos) {
			c.graphicsMod.ShowChunk(graphics.ChunkEvent{
				PositionX: pos.x,
				PositionY: pos.y,
				PositionZ: pos.z,
			})
		}
		return false
	})
	old.forEach(func(pos chunkPos) bool {
		if !new.contains(pos) {
			c.graphicsMod.HideChunk(graphics.ChunkEvent{
				PositionX: pos.x,
				PositionY: pos.y,
				PositionZ: pos.z,
			})
		}
		return false
	})
	c.lastChunkPos = newChunkPos
}
