package player_test

import (
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

func TestModuleNew(t *testing.T) {
	t.Parallel()

	t.Run("return is non-nil", func(t *testing.T) {
		t.Parallel()

		mod := player.New(world.FnModule{}, settings.FnRepository{}, &view.FnModule{})

		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})

	t.Run("panic on nil settingsMod", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()

		player.New(world.FnModule{}, nil, &view.FnModule{})
	})
	t.Run("panic on nil world mod", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()

		player.New(nil, settings.FnRepository{}, &view.FnModule{})
	})
	t.Run("panic on nil view module", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()

		player.New(world.FnModule{}, settings.FnRepository{}, nil)
	})
	t.Run("nothing is loded by default", func(t *testing.T) {
		t.Parallel()
		expected := false
		var loaded bool
		worldMod := world.FnModule{
			FnLoadChunk: func(pos chunk.ChunkCoordinate) {
				loaded = true
			},
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 0
			},
		}

		player.New(worldMod, settingsMod, &view.FnModule{})

		if loaded != expected {
			t.Fatal("expected no chunk to be loaded, but one was")
		}
	})
}

func TestModuleUpdatePlayerPosition(t *testing.T) {
	t.Parallel()
	t.Run("when player position is moved, the right chunks are loaded and unloaded", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expectedLoad := map[chunk.ChunkCoordinate]struct{}{}
		expectedUnload := map[chunk.ChunkCoordinate]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expectedLoad[chunk.ChunkCoordinate{
					X: 3,
					Y: y,
					Z: z,
				}] = struct{}{}
				expectedUnload[chunk.ChunkCoordinate{
					X: -2,
					Y: y,
					Z: z,
				}] = struct{}{}
			}
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
			FnGetChunkSize: func() uint32 {
				return chunkSize
			},
		}
		worldMod := &world.FnModule{}
		playerMod := player.New(worldMod, settingsMod, &view.FnModule{})
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actualLoaded := map[chunk.ChunkCoordinate]struct{}{}
		actualUnloaded := map[chunk.ChunkCoordinate]struct{}{}

		worldMod.FnLoadChunk = func(pos chunk.ChunkCoordinate) {
			actualLoaded[pos] = struct{}{}
		}
		worldMod.FnUnloadChunk = func(pos chunk.ChunkCoordinate) {
			actualUnloaded[pos] = struct{}{}
		}

		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: chunkSize + 5,
			Y: 0,
			Z: 0,
		})

		if !reflect.DeepEqual(expectedLoad, actualLoaded) {
			t.Fatalf("expected to load %v but got %v", expectedLoad, actualLoaded)
		}
		if !reflect.DeepEqual(expectedUnload, actualUnloaded) {
			t.Fatalf("expected to unload %v but got %v", expectedUnload, actualUnloaded)
		}
	})

	t.Run("when player position is moved diagonally, new chunks are shown", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expected := map[chunk.ChunkCoordinate]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-3); x <= 1; x++ {
				expected[chunk.ChunkCoordinate{
					X: x,
					Y: y,
					Z: -3,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-3); z <= 1; z++ {
				expected[chunk.ChunkCoordinate{
					X: -3,
					Y: y,
					Z: z,
				}] = struct{}{}
			}
		}
		worldMod := &world.FnModule{}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
			FnGetChunkSize: func() uint32 {
				return chunkSize
			},
		}
		playerMod := player.New(worldMod, settingsMod, &view.FnModule{})
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actual := map[chunk.ChunkCoordinate]struct{}{}
		worldMod.FnLoadChunk = func(pos chunk.ChunkCoordinate) {
			actual[pos] = struct{}{}
		}

		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5 - chunkSize,
			Y: 0,
			Z: 5 - chunkSize,
		})

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})

	t.Run("when player position is moved diagonally, old chunks are hidden", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expected := map[chunk.ChunkCoordinate]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-2); x <= 2; x++ {
				expected[chunk.ChunkCoordinate{
					X: x,
					Y: y,
					Z: 2,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expected[chunk.ChunkCoordinate{
					X: 2,
					Y: y,
					Z: z,
				}] = struct{}{}
			}
		}
		worldMod := &world.FnModule{}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
			FnGetChunkSize: func() uint32 {
				return chunkSize
			},
		}
		playerMod := player.New(worldMod, settingsMod, &view.FnModule{})
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 5,
		})
		actual := map[chunk.ChunkCoordinate]struct{}{}
		worldMod.FnUnloadChunk = func(pos chunk.ChunkCoordinate) {
			actual[pos] = struct{}{}
		}

		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5 - chunkSize,
			Y: 0,
			Z: 5 - chunkSize,
		})

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %v but got %v", expected, actual)
		}
	})
}

func TestUpdateViewWithoutPos(t *testing.T) {
	t.Parallel()
	viewMod := &view.FnModule{
		FnUpdateView: func(viewState view.ViewState) {
			t.Fatal("expected update view to not be called, but it was")
		},
		FnUpdateSelection: func() {
			t.Fatal("expected update selection to not be called, but it was")
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, viewMod)
	playerMod.UpdatePlayerDirection(player.DirectionEvent{})
}

func TestWorldViewUpdateWithoutDirection(t *testing.T) {
	t.Parallel()
	viewMod := &view.FnModule{
		FnUpdateView: func(viewState view.ViewState) {
			t.Fatal("expected update view to not be called, but it was")
		},
		FnUpdateSelection: func() {
			t.Fatal("expected update selection to not be called, but it was")
		},
	}
	playerMod := player.New(world.FnModule{}, settings.FnRepository{}, viewMod)
	playerMod.UpdatePlayerPosition(player.PositionEvent{})
}

func TestWorldViewUpdateWithPosAndDir(t *testing.T) {
	t.Parallel()
	updatedSelection := false
	expectedViewState := view.ViewState{
		Pos: [3]float64{1, 2, 3},
		Dir: mgl.Quat{
			W: 1,
			V: [3]float64{2, 3, 4},
		},
	}
	var actualViewState view.ViewState
	viewMod := &view.FnModule{
		FnUpdateView: func(viewState view.ViewState) {
			actualViewState = viewState
		},
		FnUpdateSelection: func() {
			updatedSelection = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, viewMod)
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: 1, Y: 2, Z: 3})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: mgl.Quat{
			W: 1,
			V: [3]float64{2, 3, 4},
		},
	})

	if actualViewState != expectedViewState {
		t.Fatalf("expected view to receive view state %v but got %v", expectedViewState, actualViewState)
	}
	if !updatedSelection {
		t.Fatal("failed to update selection")
	}
	updatedSelection = false

	expected2 := view.ViewState{
		Pos: [3]float64{7, 8, 9},
		Dir: mgl.Quat{
			W: 1,
			V: [3]float64{2, 3, 4},
		},
	}
	playerMod.UpdatePlayerPosition(player.PositionEvent{7, 8, 9})

	if actualViewState != expected2 {
		t.Fatalf("expected world to receive view state %v but got %v", expected2, actualViewState)
	}
	if !updatedSelection {
		t.Fatal("failed to update selection")
	}
}

func TestChunksLoadedOnFirstPositionUpdate(t *testing.T) {
	t.Parallel()
	expectedLoadCall := true
	actualLoadCall := false
	expectedUnloadCall := false
	actualUnloadCall := false
	worldMod := world.FnModule{
		FnLoadChunk: func(p chunk.ChunkCoordinate) {
			actualLoadCall = true
		},
		FnUnloadChunk: func(p chunk.ChunkCoordinate) {
			actualUnloadCall = true
		},
	}
	playerMod := player.New(worldMod, settings.FnRepository{}, &view.FnModule{})
	playerMod.UpdatePlayerPosition(player.PositionEvent{1, 1, 1})

	if expectedLoadCall != actualLoadCall {
		t.Fatal("expected load chunk to be called, but it wasn't")
	}
	if expectedUnloadCall != actualUnloadCall {
		t.Fatal("expected unload chunk to never be called, but it was")
	}
}

func TestViewUpdateAfterLoadingChunks(t *testing.T) {
	x := 1
	worldMod := world.FnModule{
		FnLoadChunk: func(p chunk.ChunkCoordinate) {
			x -= 1
		},
	}
	viewMod := view.FnModule{
		FnUpdateView: func(view.ViewState) {
			x *= 2
		},
	}
	expected := 0
	playerMod := player.New(worldMod, settings.FnRepository{}, &viewMod)
	playerMod.UpdatePlayerPosition(player.PositionEvent{1, 1, 1})
	if x != expected {
		t.Fatal("player updated view before loading chunks")
	}
}

func TestViewUpdateAfterUnloadingChunks(t *testing.T) {
	x := 1
	worldMod := world.FnModule{
		FnUnloadChunk: func(p chunk.ChunkCoordinate) {
			x -= 1
		},
	}
	viewMod := view.FnModule{
		FnUpdateView: func(view.ViewState) {
			x *= 2
		},
	}
	expected := 0
	playerMod := player.New(worldMod, settings.FnRepository{}, &viewMod)
	playerMod.UpdatePlayerPosition(player.PositionEvent{1, 1, 1})
	playerMod.UpdatePlayerPosition(player.PositionEvent{0, 1, 1})
	if x != expected {
		t.Fatal("player updated view before unloading chunks")
	}
}

func TestPlayerScrollDown(t *testing.T) {
	blockRemoved := false
	worldMod := world.FnModule{
		FnRemoveBlock: func(chunk.VoxelCoordinate) {
			blockRemoved = true
		},
	}
	var nodeRemoved, updated bool
	viewMod := view.FnModule{
		FnRemoveNode: func(chunk.VoxelCoordinate) {
			nodeRemoved = true
		},
		FnUpdateSelection: func() {
			updated = true
		},
		FnGetSelection: func() (chunk.VoxelCoordinate, bool) {
			return chunk.VoxelCoordinate{}, true
		},
	}

	playerMod := player.New(worldMod, settings.FnRepository{}, &viewMod)
	playerMod.UpdatePlayerAction(player.ActionEvent{
		Scroll: player.ScrollDown,
	})
	if !blockRemoved || !nodeRemoved || !updated {
		t.Fatal("failed to remove block, node, or update sel")
	}
}

func TestPlayerScrollUp(t *testing.T) {
	worldMod := world.FnModule{
		FnRemoveBlock: func(chunk.VoxelCoordinate) {
			t.Fatal("called remove block on scroll up")
		},
	}
	viewMod := view.FnModule{
		FnRemoveNode: func(chunk.VoxelCoordinate) {
			t.Fatal("called remove node on scroll up")
		},
		FnUpdateSelection: func() {
			t.Fatal("called update selection on scroll up")
		},
	}

	playerMod := player.New(worldMod, settings.FnRepository{}, &viewMod)
	playerMod.UpdatePlayerAction(player.ActionEvent{
		Scroll: player.ScrollUp,
	})
}

func withinError(x, y float64, diff float64) bool {
	if x+diff > y && x-diff < y {
		return true
	}
	return false
}

const errMargin = 0.000001
