package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/veandco/go-sdl2/sdl"
)

func (m *ParallelModule) CreateWindow(title string) error {
	done := make(chan error)
	m.do <- func() {
		done <- m.c.createWindow(title)
	}
	return <-done
}

func (m *ParallelModule) ShowWindow() {
	done := make(chan struct{})
	m.do <- func() {
		m.c.showWindow()
		close(done)
	}
	<-done
}

func (m *ParallelModule) PollEvent() (sdl.Event, bool) {
	type returns struct {
		evt sdl.Event
		ok  bool
	}
	done := make(chan returns)
	m.do <- func() {
		evt, ok := m.c.pollEvent()
		done <- returns{evt, ok}
	}
	d := <-done
	return d.evt, d.ok
}

func (m *ParallelModule) LoadChunk(chunk chunk.Chunk) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.loadChunk(chunk)
		close(done)
	}
	<-done
}

func (m *ParallelModule) UpdateChunk(chunk chunk.Chunk) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updateChunk(chunk)
		close(done)
	}
	<-done
}

func (m *ParallelModule) UnloadChunk(pos chunk.ChunkCoordinate) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.unloadChunk(pos)
		close(done)
	}
	<-done
}

func (m *ParallelModule) UpdateView(viewableChunks map[chunk.ChunkCoordinate]struct{}, viewMat mgl.Mat4) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updateView(viewableChunks, viewMat)
		close(done)
	}
	<-done
}

func (m *ParallelModule) UpdateSelection(selectedVoxel chunk.VoxelCoordinate, selected bool) {
	done := make(chan struct{})
	m.do <- func() {
		m.c.updateSelection(selectedVoxel, selected)
		close(done)
	}
	<-done
}

func (m *ParallelModule) DestroyWindow() error {
	return m.c.destroyWindow()
}

func (m *ParallelModule) Render() {
	done := make(chan struct{})
	m.do <- func() {
		m.c.render()
		close(done)
	}
	<-done
}
