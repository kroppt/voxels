package world

import (
	"container/list"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/repositories/settings"
)

type core struct {
	graphicsMod    graphics.Interface
	generator      Generator
	settingsRepo   settings.Interface
	cacheMod       cache.Interface
	viewMod        view.Interface
	loadedChunks   map[chunk.ChunkCoordinate]*chunkState
	pendingActions map[chunk.ChunkCoordinate]*list.List
}

type chunkState struct {
	ch       chunk.Chunk
	modified bool
}

func (c *core) loadChunk(pos chunk.ChunkCoordinate) {
	if _, ok := c.loadedChunks[pos]; ok {
		panic("tried to load already-loaded chunk")
	}
	ch, ok := c.cacheMod.Load(pos)
	actions := list.New()
	if !ok {
		ch, actions = c.generator.GenerateChunk(pos)
	}
	var root *view.Octree
	ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
		if ch.BlockType(vc) != chunk.BlockTypeAir {
			root = root.AddLeaf(&vc)
		}
	})
	c.viewMod.AddTree(pos, root)

	c.loadedChunks[pos] = &chunkState{
		ch:       ch,
		modified: false,
	}
	c.handlePendingActions(actions)
	if _, ok := c.pendingActions[pos]; ok {
		c.performPendingActions(pos)
	}
	c.graphicsMod.LoadChunk(ch)
}

func (c *core) unloadChunk(pos chunk.ChunkCoordinate) {
	cs, ok := c.loadedChunks[pos]
	if !ok {
		panic("tried to unload a chunk that is not loaded")
	}
	if cs.modified {
		c.cacheMod.Save(cs.ch)
	}
	c.viewMod.RemoveTree(pos)
	delete(c.loadedChunks, pos)
	c.graphicsMod.UnloadChunk(pos)
}

func (c *core) handlePendingActions(actions *list.List) {
	// go through all of the actions
	// perform any actions right now that can be performed because those chunks are loaded
	// for the ones that you can't do now, save them for later
	immediateUpdateChunks := map[chunk.ChunkCoordinate]struct{}{}
	for action := actions.Front(); action != nil; action = action.Next() {
		pa := action.Value.(chunk.PendingAction)
		if _, ok := c.loadedChunks[pa.ChPos]; ok {
			immediateUpdateChunks[pa.ChPos] = struct{}{}
		}
		if _, ok := c.pendingActions[pa.ChPos]; ok {
			c.pendingActions[pa.ChPos].PushBack(pa)
		} else {
			c.pendingActions[pa.ChPos] = list.New()
			c.pendingActions[pa.ChPos].PushBack(pa)
		}
	}

	for otherChunk := range immediateUpdateChunks {
		c.performPendingActions(otherChunk)
		c.graphicsMod.UpdateChunk(c.loadedChunks[otherChunk].ch)
	}
}

func (c *core) performPendingActions(cc chunk.ChunkCoordinate) {
	if _, ok := c.loadedChunks[cc]; !ok {
		panic("attempted to perform pending actions on a chunk that isn't loaded")
	}
	if _, ok := c.pendingActions[cc]; !ok {
		panic("attempted to perform pending actions on a chunk that doesn't have any")
	}
	actions := c.pendingActions[cc]
	for action := actions.Front(); action != nil; action = action.Next() {
		pa := action.Value.(chunk.PendingAction)
		ch := c.loadedChunks[pa.ChPos]
		if pa.HideFace {
			ch.ch.AddAdjacency(pa.VoxPos, pa.Face)
		} else {
			ch.ch.RemoveAdjacency(pa.VoxPos, pa.Face)
		}
	}
	c.loadedChunks[cc].modified = true
	delete(c.pendingActions, cc)
}

func (c *core) quit() {
	for key, actions := range c.pendingActions {
		ch, ok := c.cacheMod.Load(key)
		if !ok {
			ch, _ = c.generator.GenerateChunk(key)
		}
		for action := actions.Front(); action != nil; action = action.Next() {
			pa := action.Value.(chunk.PendingAction)
			if pa.HideFace {
				ch.AddAdjacency(pa.VoxPos, pa.Face)
			} else {
				ch.RemoveAdjacency(pa.VoxPos, pa.Face)
			}
		}
		c.cacheMod.Save(ch)
	}

	for pos, cs := range c.loadedChunks {
		if cs.modified {
			c.cacheMod.Save(cs.ch)
		}
		c.graphicsMod.UnloadChunk(pos)
	}
	c.cacheMod.Close()
}

func (c *core) countLoadedChunks() int {
	return len(c.loadedChunks)
}

func (c *core) getBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	key := chunk.VoxelCoordToChunkCoord(pos, c.settingsRepo.GetChunkSize())
	if _, ok := c.loadedChunks[key]; !ok {
		panic("tried to get block from non-loaded chunk")
	}
	return c.loadedChunks[key].ch.BlockType(pos)
}

func (c *core) removeBlock(vc chunk.VoxelCoordinate) {
	cc := chunk.VoxelCoordToChunkCoord(vc, c.settingsRepo.GetChunkSize())
	cs, ok := c.loadedChunks[cc]
	if !ok {
		panic("tried to a remove a block from a chunk that isn't loaded")
	}
	actions := cs.ch.SetBlockType(vc, chunk.BlockTypeAir)
	cs.modified = true
	c.handlePendingActions(actions)
	c.viewMod.RemoveNode(vc)
	c.graphicsMod.UpdateChunk(cs.ch)
}

func (c *core) addBlock(vc chunk.VoxelCoordinate, bt chunk.BlockType) {
	key := chunk.VoxelCoordToChunkCoord(vc, c.settingsRepo.GetChunkSize())
	cs, ok := c.loadedChunks[key]
	if !ok {
		panic("tried to add a block from a chunk that isn't loaded")
	}
	actions := cs.ch.SetBlockType(vc, bt)
	cs.modified = true
	c.handlePendingActions(actions)
	c.viewMod.AddNode(vc)
	c.graphicsMod.UpdateChunk(cs.ch)
}
