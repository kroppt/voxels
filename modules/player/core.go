package player

import (
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	worldMod     world.Interface
	settingsMod  settings.Interface
	graphicsMod  graphics.Interface
	lastChunkPos chunk.ChunkCoordinate
	posAssigned  bool
	position     PositionEvent
	dirAssigned  bool
	direction    DirectionEvent
	firstLoad    bool
}

// chunkRange is the range of chunks between Min and Max.
type chunkRange struct {
	Min chunk.ChunkCoordinate
	Max chunk.ChunkCoordinate
}

func toVoxelPos(playerPos PositionEvent) chunk.VoxelCoordinate {
	x, y, z := playerPos.X, playerPos.Y, playerPos.Z
	if x < 0 {
		x--
	}
	if y < 0 {
		y--
	}
	if z < 0 {
		z--
	}
	return chunk.VoxelCoordinate{
		X: int32(x),
		Y: int32(y),
		Z: int32(z),
	}
}

// forEach executes the given function on every position in the this ChunkRange.
// The return of fn indices whether to stop iterating
func (rng chunkRange) forEach(fn func(pos chunk.ChunkCoordinate) bool) {
	for x := rng.Min.X; x <= rng.Max.X; x++ {
		for y := rng.Min.Y; y <= rng.Max.Y; y++ {
			for z := rng.Min.Z; z <= rng.Max.Z; z++ {
				stop := fn(chunk.ChunkCoordinate{X: x, Y: y, Z: z})
				if stop {
					return
				}
			}
		}
	}
}

// contains returns whether this ChunkRange contains the given pos.
func (rng chunkRange) contains(pos chunk.ChunkCoordinate) bool {
	if pos.X < rng.Min.X || pos.X > rng.Max.X {
		return false
	}
	if pos.Y < rng.Min.Y || pos.Y > rng.Max.Y {
		return false
	}
	if pos.Z < rng.Min.Z || pos.Z > rng.Max.Z {
		return false
	}
	return true
}

func (c *core) viewState() world.ViewState {
	if !c.dirAssigned || !c.posAssigned {
		panic("direction or position not assigned, unintended use")
	}
	return world.ViewState{
		Pos: [3]float64{c.position.X, c.position.Y, c.position.Z},
		Dir: c.direction.Rotation,
	}
}

func (c *core) updatePosition(posEvent PositionEvent) {
	newChunkPos := chunk.VoxelCoordToChunkCoord(toVoxelPos(posEvent), c.settingsMod.GetChunkSize())
	renderDistance := int32(c.settingsMod.GetRenderDistance())
	old := chunkRange{
		Min: chunk.ChunkCoordinate{
			X: c.lastChunkPos.X - renderDistance,
			Y: c.lastChunkPos.Y - renderDistance,
			Z: c.lastChunkPos.Z - renderDistance,
		},
		Max: chunk.ChunkCoordinate{
			X: c.lastChunkPos.X + renderDistance,
			Y: c.lastChunkPos.Y + renderDistance,
			Z: c.lastChunkPos.Z + renderDistance,
		},
	}
	new := chunkRange{
		Min: chunk.ChunkCoordinate{
			X: newChunkPos.X - renderDistance,
			Y: newChunkPos.Y - renderDistance,
			Z: newChunkPos.Z - renderDistance,
		},
		Max: chunk.ChunkCoordinate{
			X: newChunkPos.X + renderDistance,
			Y: newChunkPos.Y + renderDistance,
			Z: newChunkPos.Z + renderDistance,
		},
	}
	new.forEach(func(pos chunk.ChunkCoordinate) bool {
		if !old.contains(pos) || c.firstLoad {
			c.worldMod.LoadChunk(pos)
		}
		return false
	})
	old.forEach(func(pos chunk.ChunkCoordinate) bool {
		if !new.contains(pos) && !c.firstLoad {
			c.worldMod.UnloadChunk(pos)
		}
		return false
	})
	if c.firstLoad {
		c.firstLoad = false
	}
	c.lastChunkPos = newChunkPos

	c.posAssigned = true
	c.position = posEvent
	if c.dirAssigned {
		c.worldMod.UpdateView(c.viewState())
	}
}

func (c *core) updateDirection(dirEvent DirectionEvent) {
	c.dirAssigned = true
	c.direction = dirEvent
	if c.posAssigned {
		c.worldMod.UpdateView(c.viewState())
	}
}
