package view_test

import (
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/repositories/settings"
)

func TestRequiredSubModules(t *testing.T) {
	t.Parallel()
	t.Run("graphics module is not nil", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		view.New(nil, settings.FnRepository{})
	})
	t.Run("settings repo is not nil", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		view.New(graphics.FnModule{}, nil)
	})
}

func TestUpdateSelectionCallsGraphics(t *testing.T) {
	graphicsMod := graphics.FnModule{
		FnUpdateSelection: func(vc chunk.VoxelCoordinate, b bool) {
			if b {
				t.Fatal("should not have received a selection")
			}
		},
	}
	viewMod := view.New(graphicsMod, settings.FnRepository{
		FnGetChunkSize: func() uint32 { return 1 },
	})
	viewMod.UpdateSelection()
}

func TestCannotGetSelectionWithoutViewState(t *testing.T) {
	expected := false
	viewMod := view.New(graphics.FnModule{}, settings.FnRepository{})
	_, actual := viewMod.GetSelection()
	if actual != expected {
		t.Fatal("expected to get a false selection, but got a true one")
	}
}

func TestGetSelectionValidEmptyTree(t *testing.T) {
	expected := false
	viewMod := view.New(graphics.FnModule{}, settings.FnRepository{})
	viewMod.UpdateView(view.ViewState{
		Pos: [3]float64{},
		Dir: mgl.QuatIdent(),
	})
	_, actual := viewMod.GetSelection()
	if actual != expected {
		t.Fatal("expected to get a false selection, but got a true one")
	}
}

func TestGetSelectionValid(t *testing.T) {
	expectedSelection := chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1}
	viewMod := view.New(graphics.FnModule{}, settings.FnRepository{})
	viewMod.UpdateView(view.ViewState{
		Pos: [3]float64{0.5, 0.5, 0.5},
		Dir: mgl.QuatIdent(),
	})
	viewMod.AddNode(chunk.VoxelCoordinate{X: 0, Y: 0, Z: -1})
	actualSelection, selected := viewMod.GetSelection()
	if !selected {
		t.Fatal("expected to get a true selection, but got a false one")
	}
	if actualSelection != expectedSelection {
		t.Fatalf("expected to select %v but got %v\n", expectedSelection, actualSelection)
	}
}
func TestUpdateViewCallsGraphics(t *testing.T) {
	expected := true
	actual := false
	graphicsMod := graphics.FnModule{
		FnUpdateView: func(m1 map[chunk.ChunkCoordinate]struct{}, m2 mgl.Mat4) {
			actual = true
		},
	}
	viewMod := view.New(graphicsMod, settings.FnRepository{})
	viewMod.UpdateView(view.ViewState{})
	if actual != expected {
		t.Fatal("expected view to update view in graphics, but it did not")
	}
}

// TODO test accuracy of view update, calculate by hand

func TestFrustumCulling(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc       string
		expect     map[chunk.ChunkCoordinate]struct{}
		viewState  view.ViewState
		fov        float64
		far        float64
		near       float64
		width      uint32
		height     uint32
		renderDist uint32
		chunkSize  uint32
	}{
		{
			desc: "simple frustum culling",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: 0, Z: 0}:  {},
				{X: 0, Y: 0, Z: -1}: {},
			},
			viewState: view.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        33.398488467987,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  1,
		},
		{
			desc: "frustum culling large chunks",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: 0, Z: 0}:    {},
				{X: 0, Y: 0, Z: -1}:   {},
				{X: -1, Y: 0, Z: -1}:  {},
				{X: 0, Y: -1, Z: -1}:  {},
				{X: -1, Y: -1, Z: -1}: {},
			},
			viewState: view.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        70,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  3,
		},
		{
			desc: "frustum culling wide angle",
			expect: map[chunk.ChunkCoordinate]struct{}{
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
			},
			viewState: view.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.QuatIdent(),
			},
			fov:        89.5,
			far:        10,
			near:       0.1,
			width:      1,
			height:     1,
			renderDist: 1,
			chunkSize:  1,
		},
		{
			desc: "edge case found in production",
			expect: map[chunk.ChunkCoordinate]struct{}{
				{X: 0, Y: -1, Z: -1}: {},
				{X: 0, Y: 0, Z: -1}:  {},
				{X: 0, Y: 1, Z: -1}:  {},
				{X: 1, Y: -1, Z: -1}: {},
				{X: 1, Y: 0, Z: -1}:  {},
				{X: 1, Y: 1, Z: -1}:  {},
				{X: 1, Y: -1, Z: 0}:  {},
				{X: 1, Y: 0, Z: 0}:   {},
				{X: 1, Y: 1, Z: 0}:   {},
				{X: 0, Y: 0, Z: 0}:   {},
			},
			viewState: view.ViewState{
				Pos: [3]float64{0.5, 0.5, 0.5},
				Dir: mgl.Quat{
					W: 0.9238795325112867,
					V: mgl.Vec3{0, -0.3826834323650898, 0},
				},
			},
			fov:        60,
			far:        100,
			near:       0.1,
			width:      1280,
			height:     720,
			renderDist: 1,
			chunkSize:  1,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			actualViewedChunks := map[chunk.ChunkCoordinate]struct{}{}
			graphicsMod := &graphics.FnModule{
				FnUpdateView: func(viewChunks map[chunk.ChunkCoordinate]struct{}, _ mgl.Mat4) {
					actualViewedChunks = viewChunks
				},
			}
			settingsMod := settings.FnRepository{
				FnGetFOV: func() float64 {
					return tC.fov
				},
				FnGetFar: func() float64 {
					return tC.far
				},
				FnGetNear: func() float64 {
					return tC.near
				},
				FnGetResolution: func() (uint32, uint32) {
					return tC.width, tC.height
				},
				FnGetRenderDistance: func() uint32 {
					return tC.renderDist
				},
				FnGetChunkSize: func() uint32 {
					return tC.chunkSize
				},
			}
			viewMod := view.New(graphicsMod, settingsMod)
			viewMod.UpdateView(view.ViewState{
				Pos: tC.viewState.Pos,
				Dir: tC.viewState.Dir,
			})

			if !reflect.DeepEqual(tC.expect, actualViewedChunks) {
				t.Fatalf("expected viewed chunks: %v but got viewed chunks %v", tC.expect, actualViewedChunks)
			}
		})
	}
}

func BenchmarkViewAddRemoveNode(b *testing.B) {
	voxCoord := chunk.VoxelCoordinate{
		X: 0,
		Y: 0,
		Z: 0,
	}
	viewMod := view.New(graphics.FnModule{}, &settings.Repository{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		viewMod.AddNode(voxCoord)
		viewMod.RemoveNode(voxCoord)
	}
}

func BenchmarkViewAddRemoveChunk(b *testing.B) {
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{X: 0, Y: 0, Z: 0}, 5)
	viewMod := view.New(graphics.FnModule{}, &settings.Repository{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
			viewMod.AddNode(vc)
		})
		ch.ForEachVoxel(func(vc chunk.VoxelCoordinate) {
			viewMod.RemoveNode(vc)
		})
	}
}
