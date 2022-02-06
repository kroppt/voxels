package world

import "github.com/kroppt/voxels/chunk"

func (m *ParallelModule) LoadChunk(pos chunk.ChunkCoordinate) {
	m.do <- func() {
		m.c.loadChunk(pos)
	}
}

func (m *ParallelModule) UnloadChunk(pos chunk.ChunkCoordinate) {
	m.do <- func() {
		m.c.unloadChunk(pos)
	}
}

func (m *ParallelModule) Quit() {
	m.c.quit()
}

func (m *ParallelModule) CountLoadedChunks() int {
	done := make(chan int)
	m.do <- func() {
		done <- m.c.countLoadedChunks()
	}
	return <-done
}

func (m *ParallelModule) GetBlockType(pos chunk.VoxelCoordinate) chunk.BlockType {
	done := make(chan chunk.BlockType)
	m.do <- func() {
		done <- m.c.getBlockType(pos)
	}
	return <-done
}

func (m *ParallelModule) RemoveBlock(vc chunk.VoxelCoordinate) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.removeBlock(vc)
		close(done)
	}
	<-done
}

func (m *ParallelModule) AddBlock(vc chunk.VoxelCoordinate, bt chunk.BlockType) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.addBlock(vc, bt)
		close(done)
	}
	<-done
}
