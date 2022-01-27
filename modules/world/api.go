package world

// world is the intermediary between the player and graphics
// the world should contain the true state of chunks/voxels
// for purposes of editing the world and performing lighting calcs
// beneath the API, the world module will send the graphics module
// chunk updates in the form of flat data for it to continuously render

type Interface interface {
	// LoadChunk will generate a new chunk or load one from a file
	// the implementation will store it internally as flat data.
	// Also, this will initially inform the graphics mod. of the chunk
	// so it can be rendered.
	LoadChunk(ChunkEvent)
	// UnloadChunk will remove a chunk from memory, saving it to
	// a file if necessary. It will also tell the graphics mod. that
	// this chunk should be removed from its memory and deallocated in opengl
	UnloadChunk(ChunkEvent)
	// SaveLoadedChunks's purpose is to save the state of the world that is left
	// modified in memory in the case where the program is shutting down/explicitly
	// told to save.
	// SaveLoadedChunks()
	// SetVoxelAt sets the attributes of the voxel specified at the given position.
	// This will change the flat data of a chunk, and thus will need to be communicated
	// to the graphics module.
	// RemoveVoxel is not necessary with the current planned implementation where
	// the flat data per chunk is a fixed size, but this could change in the future
	// if we figure out how to "not store air".
	// SetVoxelAt(VoxelEvent, VoxelAttribs)
	// UpdateFrustumCulling is called whenever there is a change in camera position/angle
	// Since the world knows about all the chunks, and we want to keep frustum
	// culling processing out of the graphics thread, the world will accept
	// camera movement events to re-calculate frustum culling, and then tell the
	// graphics module that some chunks should be temporarily hidden, or show those
	// that have come into view again.
	// (not removed from memory/opengl allocation - just updated metadata that will
	// cause certain chunks to be rendered on the next render cycle)
	// UpdateFrustumCulling(CameraUpdateEvent)
}

type ChunkEvent struct {
	PositionX int32
	PositionY int32
	PositionZ int32
}

// CameraUpdateEvent describes the position and direction of the camera.
// type CameraUpdateEvent struct {
// 	X   int32
// 	Y   int32
// 	Z   int32
// 	Rot mgl.Quat
// }

// type VoxelAttribs struct {
// 	BlockType int32
// }

// VoxelEvent contains voxel position information.
// type VoxelEvent struct {
// 	X int32
// 	Y int32
// 	Z int32
// }

func (m *Module) LoadChunk(ChunkEvent) {

}
func (m *Module) UnloadChunk(ChunkEvent) {

}

// func (m *Module) SaveLoadedChunks() {

// }

// func (m *Module) SetVoxelAt(VoxelEvent, VoxelAttribs) {

// }

// func (m *Module) UpdateFrustumCulling(CameraUpdateEvent) {

// }

type FnModule struct {
	FnLoadChunk   func(ChunkEvent)
	FnUnloadChunk func(ChunkEvent)
	// FnSaveLoadedChunks func()
	// FnSetVoxelAt       func(VoxelEvent, VoxelAttribs)
	// FnUpdateFrustumCulling func(CameraUpdateEvent)
}

func (fn *FnModule) LoadChunk(evt ChunkEvent) {
	if fn.FnLoadChunk != nil {
		fn.FnLoadChunk(evt)
	}
}
func (fn *FnModule) UnloadChunk(evt ChunkEvent) {
	if fn.FnUnloadChunk != nil {
		fn.FnUnloadChunk(evt)
	}
}

// func (fn *FnModule) SaveLoadedChunks() {
// 	if fn.FnSaveLoadedChunks != nil {
// 		fn.FnSaveLoadedChunks()
// 	}
// }

// func (fn *FnModule) SetVoxelAt(evt VoxelEvent, attrs VoxelAttribs) {
// 	if fn.FnSetVoxelAt != nil {
// 		fn.SetVoxelAt(evt, attrs)
// 	}
// }

// func (fn *FnModule) UpdateFrustumCulling(evt CameraUpdateEvent) {
// 	if fn.FnUpdateFrustumCulling != nil {
// 		fn.UpdateFrustumCulling(evt)
// 	}
// }
