package player_test

import (
	"reflect"
	"testing"

	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

func TestModuleNew(t *testing.T) {
	t.Parallel()

	t.Run("return is non-nil", func(t *testing.T) {
		t.Parallel()
		worldMod := &world.FnModule{}
		settingsMod := settings.FnRepository{}

		mod := player.New(worldMod, settingsMod, 1)

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
		worldMod := &world.FnModule{}

		player.New(worldMod, nil, 1)
	})

	t.Run("nothing is loded by default", func(t *testing.T) {
		t.Parallel()
		expected := false
		var loaded bool
		worldMod := &world.FnModule{
			FnLoadChunk: func(chunkEvent world.ChunkEvent) {
				loaded = true
			},
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 0
			},
		}

		player.New(worldMod, settingsMod, 1)

		if loaded != expected {
			t.Fatal("expected no chunk to be loaded, but one was")
		}
	})
}

func TestModuleUpdatePlayerPosition(t *testing.T) {
	t.Run("when player position is moved, the right chunks are loaded and unloaded", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expectedLoad := map[world.ChunkEvent]struct{}{}
		expectedUnload := map[world.ChunkEvent]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expectedLoad[world.ChunkEvent{
					PositionX: 3,
					PositionY: y,
					PositionZ: z,
				}] = struct{}{}
				expectedUnload[world.ChunkEvent{
					PositionX: -2,
					PositionY: y,
					PositionZ: z,
				}] = struct{}{}
			}
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
		}
		worldMod := &world.FnModule{}
		playerMod := player.New(worldMod, settingsMod, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actualLoaded := map[world.ChunkEvent]struct{}{}
		actualUnloaded := map[world.ChunkEvent]struct{}{}

		worldMod.FnLoadChunk = func(chunkEvent world.ChunkEvent) {
			actualLoaded[chunkEvent] = struct{}{}
		}
		worldMod.FnUnloadChunk = func(chunkEvent world.ChunkEvent) {
			actualUnloaded[chunkEvent] = struct{}{}
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
		expected := map[world.ChunkEvent]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-3); x <= 1; x++ {
				expected[world.ChunkEvent{
					PositionX: x,
					PositionY: y,
					PositionZ: -3,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-3); z <= 1; z++ {
				expected[world.ChunkEvent{
					PositionX: -3,
					PositionY: y,
					PositionZ: z,
				}] = struct{}{}
			}
		}
		worldMod := &world.FnModule{}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
		}
		playerMod := player.New(worldMod, settingsMod, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actual := map[world.ChunkEvent]struct{}{}
		worldMod.FnLoadChunk = func(chunkEvent world.ChunkEvent) {
			actual[chunkEvent] = struct{}{}
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
		expected := map[world.ChunkEvent]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-2); x <= 2; x++ {
				expected[world.ChunkEvent{
					PositionX: x,
					PositionY: y,
					PositionZ: 2,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expected[world.ChunkEvent{
					PositionX: 2,
					PositionY: y,
					PositionZ: z,
				}] = struct{}{}
			}
		}
		worldMod := &world.FnModule{}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 2
			},
		}
		playerMod := player.New(worldMod, settingsMod, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 5,
		})
		actual := map[world.ChunkEvent]struct{}{}
		worldMod.FnUnloadChunk = func(chunkEvent world.ChunkEvent) {
			actual[chunkEvent] = struct{}{}
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

func TestNoCullingWithoutPos(t *testing.T) {
	t.Parallel()
	expected := false
	var calledUpdateView bool
	worldMod := &world.FnModule{
		FnUpdateView: func() {
			calledUpdateView = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(worldMod, settingsMod, 1)
	playerMod.UpdatePlayerDirection(player.DirectionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to not be called, but it was")
	}
}

func TestCullingWithPos(t *testing.T) {
	t.Parallel()
	expected := true
	var calledUpdateView bool
	worldMod := &world.FnModule{
		FnUpdateView: func() {
			calledUpdateView = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(worldMod, settingsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to be called, but it was not")
	}
}
