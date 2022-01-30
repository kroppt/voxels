package player_test

import (
	"math"
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
)

func TestModuleNew(t *testing.T) {
	t.Parallel()

	t.Run("return is non-nil", func(t *testing.T) {
		t.Parallel()

		mod := player.New(world.FnModule{}, settings.FnRepository{}, graphics.FnModule{}, 1)

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

		player.New(world.FnModule{}, nil, graphics.FnModule{}, 1)
	})

	t.Run("nothing is loded by default", func(t *testing.T) {
		t.Parallel()
		expected := false
		var loaded bool
		worldMod := world.FnModule{
			FnLoadChunk: func(pos chunk.Position) {
				loaded = true
			},
		}
		settingsMod := settings.FnRepository{
			FnGetRenderDistance: func() uint32 {
				return 0
			},
		}

		player.New(worldMod, settingsMod, graphics.FnModule{}, 1)

		if loaded != expected {
			t.Fatal("expected no chunk to be loaded, but one was")
		}
	})
}

func TestModuleUpdatePlayerPosition(t *testing.T) {
	t.Run("when player position is moved, the right chunks are loaded and unloaded", func(t *testing.T) {
		t.Parallel()
		const chunkSize = 10
		expectedLoad := map[chunk.Position]struct{}{}
		expectedUnload := map[chunk.Position]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expectedLoad[chunk.Position{
					X: 3,
					Y: y,
					Z: z,
				}] = struct{}{}
				expectedUnload[chunk.Position{
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
		}
		worldMod := &world.FnModule{}
		playerMod := player.New(worldMod, settingsMod, graphics.FnModule{}, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actualLoaded := map[chunk.Position]struct{}{}
		actualUnloaded := map[chunk.Position]struct{}{}

		worldMod.FnLoadChunk = func(pos chunk.Position) {
			actualLoaded[pos] = struct{}{}
		}
		worldMod.FnUnloadChunk = func(pos chunk.Position) {
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
		expected := map[chunk.Position]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-3); x <= 1; x++ {
				expected[chunk.Position{
					X: x,
					Y: y,
					Z: -3,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-3); z <= 1; z++ {
				expected[chunk.Position{
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
		}
		playerMod := player.New(worldMod, settingsMod, graphics.FnModule{}, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 0,
		})
		actual := map[chunk.Position]struct{}{}
		worldMod.FnLoadChunk = func(pos chunk.Position) {
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
		expected := map[chunk.Position]struct{}{}
		for y := int32(-2); y <= 2; y++ {
			for x := int32(-2); x <= 2; x++ {
				expected[chunk.Position{
					X: x,
					Y: y,
					Z: 2,
				}] = struct{}{}
			}
		}
		for y := int32(-2); y <= 2; y++ {
			for z := int32(-2); z <= 2; z++ {
				expected[chunk.Position{
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
		}
		playerMod := player.New(worldMod, settingsMod, graphics.FnModule{}, chunkSize)
		playerMod.UpdatePlayerPosition(player.PositionEvent{
			X: 5,
			Y: 0,
			Z: 5,
		})
		actual := map[chunk.Position]struct{}{}
		worldMod.FnUnloadChunk = func(pos chunk.Position) {
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

func TestNoCullingWithoutPos(t *testing.T) {
	t.Parallel()
	expected := false
	var calledUpdateView bool
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			calledUpdateView = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerDirection(player.DirectionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to not be called, but it was")
	}
}

func TestNoCullingWithoutDirection(t *testing.T) {
	t.Parallel()
	expected := false
	var calledUpdateView bool
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			calledUpdateView = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to not be called, but it was")
	}
}

func TestCullingWithPosAndDir(t *testing.T) {
	t.Parallel()
	expected := true
	var calledUpdateView bool
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			calledUpdateView = true
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to be called after updating direction, but it was not")
	}

	calledUpdateView = false
	playerMod.UpdatePlayerPosition(player.PositionEvent{})

	if calledUpdateView != expected {
		t.Fatal("expected update view to be called after updating position, but it was not")
	}
}

func TestFrustumCulling(t *testing.T) {
	t.Parallel()
	expectedViewedChunks := map[chunk.Position]struct{}{
		{X: 0, Y: 0, Z: 0}:  {},
		{X: 0, Y: 0, Z: -1}: {},
	}
	actualViewedChunks := map[chunk.Position]struct{}{}
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actualViewedChunks = viewChunks
		},
	}
	settingsMod := settings.New()
	settingsMod.SetFOV(33.398488467987)
	settingsMod.SetNear(0.1)
	settingsMod.SetFar(10)
	settingsMod.SetRenderDistance(1)
	settingsMod.SetResolution(1, 1)
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{0.5, 0.5, 0.5})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: mgl.QuatIdent(),
	})

	if !reflect.DeepEqual(expectedViewedChunks, actualViewedChunks) {
		t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", expectedViewedChunks, actualViewedChunks)
	}
}

func TestFrustumCullingWideAngle(t *testing.T) {
	t.Parallel()
	expectedViewedChunks := map[chunk.Position]struct{}{
		{X: 0, Y: 0, Z: 0}:    {},
		{X: 0, Y: 0, Z: -1}:   {},
		{X: -1, Y: 0, Z: -1}:  {},
		{X: -1, Y: 1, Z: -1}:  {},
		{X: -1, Y: -1, Z: -1}: {},
		{X: 0, Y: -1, Z: -1}:  {},
		{X: 0, Y: 1, Z: -1}:   {},
		{X: 1, Y: 0, Z: -1}:   {},
		{X: 1, Y: -1, Z: -1}:  {},
		{X: 1, Y: 1, Z: -1}:   {},
	}
	actualViewedChunks := map[chunk.Position]struct{}{}
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actualViewedChunks = viewChunks
		},
	}
	settingsMod := settings.New()
	settingsMod.SetFOV(89.5)
	settingsMod.SetNear(0.1)
	settingsMod.SetFar(10)
	settingsMod.SetRenderDistance(1)
	settingsMod.SetResolution(1, 1)
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{0.5, 0.5, 0.5})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: mgl.QuatIdent(),
	})

	if !reflect.DeepEqual(expectedViewedChunks, actualViewedChunks) {
		t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", expectedViewedChunks, actualViewedChunks)
	}
}

func TestFrustumCullingLargeChunks(t *testing.T) {
	t.Parallel()
	expectedViewedChunks := map[chunk.Position]struct{}{
		{X: 0, Y: 0, Z: 0}:    {},
		{X: 0, Y: 0, Z: -1}:   {},
		{X: -1, Y: 0, Z: -1}:  {},
		{X: 0, Y: -1, Z: -1}:  {},
		{X: -1, Y: -1, Z: -1}: {},
	}
	actualViewedChunks := map[chunk.Position]struct{}{}
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actualViewedChunks = viewChunks
		},
	}
	settingsMod := settings.New()
	settingsMod.SetFOV(70)
	settingsMod.SetNear(0.1)
	settingsMod.SetFar(10)
	settingsMod.SetRenderDistance(1)
	settingsMod.SetResolution(1, 1)
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 3)
	playerMod.UpdatePlayerPosition(player.PositionEvent{0.5, 0.5, 0.5})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: mgl.QuatIdent(),
	})

	if !reflect.DeepEqual(expectedViewedChunks, actualViewedChunks) {
		t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", expectedViewedChunks, actualViewedChunks)
	}
}

func TestFrustumCullingDueToPositionChange(t *testing.T) {
	t.Parallel()
	expectedViewedChunks := map[chunk.Position]struct{}{
		{X: 0, Y: 0, Z: 0}:  {},
		{X: 0, Y: 0, Z: -1}: {},
	}
	actualViewedChunks := map[chunk.Position]struct{}{}
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actualViewedChunks = viewChunks
		},
	}
	settingsMod := settings.New()
	settingsMod.SetFOV(33.398488467987)
	settingsMod.SetNear(0.1)
	settingsMod.SetFar(10)
	settingsMod.SetRenderDistance(1)
	settingsMod.SetResolution(1, 1)
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	// setting direction first without position set should not trigger a view update
	playerMod.UpdatePlayerDirection(player.DirectionEvent{
		Rotation: mgl.QuatIdent(),
	})
	if len(actualViewedChunks) != 0 {
		t.Fatal("expected update view map to be empty, but it had elements already")
	}
	playerMod.UpdatePlayerPosition(player.PositionEvent{0.5, 0.5, 0.5})

	if !reflect.DeepEqual(expectedViewedChunks, actualViewedChunks) {
		t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", expectedViewedChunks, actualViewedChunks)
	}
}

func TestViewMatrixCalculationOnDirTrigger(t *testing.T) {
	t.Parallel()
	pos := mgl.Vec3{0.5, -1, 2}
	rot := mgl.QuatIdent()
	posNeg := pos.Mul(-1)
	posMat := mgl.Translate3D(posNeg.X(), posNeg.Y(), posNeg.Z())
	expected := mgl.Ident4().Mul4(rot.Mat4()).Mul4(posMat)
	var actual mgl.Mat4
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actual = viewMat
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: pos.X(), Y: pos.Y(), Z: pos.Z()})
	playerMod.UpdatePlayerDirection(player.DirectionEvent{Rotation: rot})

	if actual != expected {
		t.Fatalf("expected graphics to receive view matrix:\n%v but got:\n%v", expected, actual)
	}
}

func TestViewMatrixCalculationOnPosTrigger(t *testing.T) {
	t.Parallel()
	pos := mgl.Vec3{0.5, -1, 2}
	rot := mgl.QuatIdent()
	posNeg := pos.Mul(-1)
	posMat := mgl.Translate3D(posNeg.X(), posNeg.Y(), posNeg.Z())
	expected := mgl.Ident4().Mul4(rot.Mat4()).Mul4(posMat)
	var actual mgl.Mat4
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actual = viewMat
		},
	}
	settingsMod := settings.FnRepository{}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerDirection(player.DirectionEvent{Rotation: rot})
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: pos.X(), Y: pos.Y(), Z: pos.Z()})

	if actual != expected {
		t.Fatalf("expected graphics to receive view matrix:\n%v but got:\n%v", expected, actual)
	}
}

func TestProjectionMatrixOnUpdateView(t *testing.T) {
	t.Parallel()
	fovRad := mgl.DegToRad(60)
	nmf, f := 1/(0.1-100), 1./math.Tan(fovRad/2.0)
	aspect := 16.0 / 9.0
	expected := mgl.Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (0.1 + 100) * nmf, -1,
		0, 0, (2. * 100 * 0.1) * nmf, 0,
	}
	var actual mgl.Mat4
	graphicsMod := &graphics.FnModule{
		FnUpdateView: func(viewChunks map[chunk.Position]struct{}, viewMat mgl.Mat4, projMat mgl.Mat4) {
			actual = projMat
		},
	}
	settingsMod := settings.FnRepository{
		FnGetFOV: func() float64 {
			return 60
		},
		FnGetFar: func() float64 {
			return 100
		},
		FnGetNear: func() float64 {
			return 0.1
		},
		FnGetResolution: func() (uint32, uint32) {
			return 1280, 720
		},
	}
	playerMod := player.New(world.FnModule{}, settingsMod, graphicsMod, 1)
	playerMod.UpdatePlayerDirection(player.DirectionEvent{Rotation: mgl.QuatIdent()})
	playerMod.UpdatePlayerPosition(player.PositionEvent{X: 0, Y: 0, Z: 0})

	if actual != expected {
		t.Fatalf("expected graphics to receive proj matrix:\n%v but got:\n%v", expected, actual)
	}
}

func TestChunksLoadedOnFirstPositionUpdate(t *testing.T) {
	t.Parallel()
	expectedLoadCall := true
	actualLoadCall := false
	expectedUnloadCall := false
	actualUnloadCall := false
	worldMod := world.FnModule{
		FnLoadChunk: func(p chunk.Position) {
			actualLoadCall = true
		},
		FnUnloadChunk: func(p chunk.Position) {
			actualUnloadCall = true
		},
	}
	playerMod := player.New(worldMod, settings.FnRepository{}, nil, 1)
	playerMod.UpdatePlayerPosition(player.PositionEvent{1, 1, 1})

	if expectedLoadCall != actualLoadCall {
		t.Fatal("expected load chunk to be called, but it wasn't")
	}
	if expectedUnloadCall != actualUnloadCall {
		t.Fatal("expected unload chunk to never be called, but it was")
	}
}
