package chunk_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kroppt/voxels/chunk"
)

func TestInvalidChunk(t *testing.T) {
	t.Parallel()
	t.Run("cannot create chunk with size 0", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		chunk.NewEmpty(chunk.Position{}, 0)
	})
	t.Run("cannot set block type out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetBlockType(chunk.VoxelCoordinate{2, 2, 2}, chunk.BlockTypeAir)
	})
	t.Run("cannot set block type out of chunk bounds complex coords", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{-1, -1, -1}, 2)
		c.SetBlockType(chunk.VoxelCoordinate{-1, -1, 1}, chunk.BlockTypeAir)
	})
	t.Run("cannot get block type out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.BlockType(chunk.VoxelCoordinate{2, 2, 2})
	})
	t.Run("cannot set lighting out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{2, 2, 2}, chunk.LightLeft, 5)
	})
	t.Run("cannot get lighting out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.Lighting(chunk.VoxelCoordinate{2, 2, 2}, chunk.LightFront)
	})
	t.Run("cannot set adjacency out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetAdjacency(chunk.VoxelCoordinate{2, 2, 2}, chunk.AdjacentBack)
	})
	t.Run("cannot get adjacency out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.Adjacency(chunk.VoxelCoordinate{2, 2, 2})
	})
	t.Run("cannot set invalid adjacency", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetAdjacency(chunk.VoxelCoordinate{0, 0, 0}, 0b1000000)
	})
	t.Run("cannot set invalid light intensity", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{0, 0, 0}, chunk.LightBack, 16)
	})
	t.Run("cannot specifcy invalid light face when setting", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{0, 0, 0}, 21, 5)
	})
	t.Run("cannot specifcy invalid light face when getting", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewEmpty(chunk.Position{0, 0, 0}, 2)
		c.Lighting(chunk.VoxelCoordinate{0, 0, 0}, 21)
	})
}

func TestChunk(t *testing.T) {
	t.Parallel()
	t.Run("test get chunk position", func(t *testing.T) {
		t.Parallel()
		expected := chunk.Position{1, 2, 3}
		chunk := chunk.NewEmpty(expected, 10)
		actual := chunk.Position()
		if actual != expected {
			t.Fatalf("expected to get chunk position %v but got %v", expected, actual)
		}
	})
	t.Run("test get chunk size", func(t *testing.T) {
		t.Parallel()
		expected := uint32(10)
		chunk := chunk.NewEmpty(chunk.Position{1, 2, 3}, expected)
		actual := chunk.Size()
		if actual != expected {
			t.Fatalf("expected to get chunk size %v but got %v", expected, actual)
		}
	})
	t.Run("check that default block type is air", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeAir
		voxelCoordinate := chunk.VoxelCoordinate{4, 4, 4}
		chunk := chunk.NewEmpty(chunk.Position{0, 0, 0}, 6)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("set block type of one voxel to dirt and get it back", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeDirt
		voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
		chunk := chunk.NewEmpty(chunk.Position{0, 0, 0}, 10)
		chunk.SetBlockType(voxelCoordinate, expected)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("set block type of one voxel to dirt and get it back offset chunk", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeDirt
		voxelCoordinate := chunk.VoxelCoordinate{12, 12, 12}
		chunk := chunk.NewEmpty(chunk.Position{1, 1, 1}, 10)
		chunk.SetBlockType(voxelCoordinate, expected)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("set block type of one voxel to dirt and get it back in negative", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeDirt
		voxelCoordinate := chunk.VoxelCoordinate{-2, -3, -4}
		chunk := chunk.NewEmpty(chunk.Position{-1, -1, -1}, 10)
		chunk.SetBlockType(voxelCoordinate, expected)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("flat data contains correct indices of voxels", func(t *testing.T) {
		t.Parallel()
		chPos := chunk.Position{0, 0, 0}
		size := int32(2)
		expectedFlatData := []float32{
			0, 0, 0, 0, 0,
			1, 0, 0, 0, 0,
			0, 1, 0, 0, 0,
			1, 1, 0, 0, 0,
			0, 0, 1, 0, 0,
			1, 0, 1, 0, 0,
			0, 1, 1, 0, 0,
			1, 1, 1, 0, 0,
		}
		ch := chunk.NewEmpty(chPos, uint32(size))
		actualFlatData := ch.GetFlatData()
		if !reflect.DeepEqual(expectedFlatData, actualFlatData) {
			t.Fatalf("expected flat data to be %v but got %v", expectedFlatData, actualFlatData)
		}
	})
}

func TestChunkVoxelAdjacencies(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		adjMask chunk.AdjacentMask
		expect  chunk.AdjacentMask
	}{
		{
			desc:    "set front adj",
			adjMask: chunk.AdjacentFront,
			expect:  chunk.AdjacentFront,
		},
		{
			desc:    "set back adj",
			adjMask: chunk.AdjacentBack,
			expect:  chunk.AdjacentBack,
		},
		{
			desc:    "set left adj",
			adjMask: chunk.AdjacentLeft,
			expect:  chunk.AdjacentLeft,
		},
		{
			desc:    "set right adj",
			adjMask: chunk.AdjacentRight,
			expect:  chunk.AdjacentRight,
		},
		{
			desc:    "set top adj",
			adjMask: chunk.AdjacentTop,
			expect:  chunk.AdjacentTop,
		},
		{
			desc:    "set bottom adj",
			adjMask: chunk.AdjacentBottom,
			expect:  chunk.AdjacentBottom,
		},
		{
			desc:    "set X adj",
			adjMask: chunk.AdjacentX,
			expect:  chunk.AdjacentX,
		},
		{
			desc:    "set Y adj",
			adjMask: chunk.AdjacentY,
			expect:  chunk.AdjacentY,
		},
		{
			desc:    "set Z adj",
			adjMask: chunk.AdjacentZ,
			expect:  chunk.AdjacentZ,
		},
		{
			desc:    "set all adj",
			adjMask: chunk.AdjacentAll,
			expect:  chunk.AdjacentAll,
		},
		{
			desc:    "set none adj",
			adjMask: chunk.AdjacentNone,
			expect:  chunk.AdjacentNone,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
			chunk := chunk.NewEmpty(chunk.Position{0, 0, 0}, 10)
			chunk.SetAdjacency(voxelCoordinate, tC.adjMask)
			actual := chunk.Adjacency(voxelCoordinate)
			if actual != tC.expect {
				t.Fatalf("expected to get adj mask %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestChunkVoxelLight(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc   string
		face   chunk.LightFace
		expect uint32
	}{
		{
			desc:   "set front light to 5",
			face:   chunk.LightFront,
			expect: 5,
		},
		{
			desc:   "set back light to 0",
			face:   chunk.LightFront,
			expect: 0,
		},
		{
			desc:   "set left light to 15",
			face:   chunk.LightLeft,
			expect: 15,
		},
		{
			desc:   "set right light to 8",
			face:   chunk.LightRight,
			expect: 8,
		},
		{
			desc:   "set bottom light to 6",
			face:   chunk.LightBottom,
			expect: 6,
		},
		{
			desc:   "set top light to 2",
			face:   chunk.LightTop,
			expect: 2,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
			chunk := chunk.NewEmpty(chunk.Position{0, 0, 0}, 10)
			chunk.SetLighting(voxelCoordinate, tC.face, tC.expect)
			actual := chunk.Lighting(voxelCoordinate, tC.face)
			if actual != tC.expect {
				t.Fatalf("expected to get light value %v but got %v", tC.expect, actual)
			}
		})
	}
}

func TestVoxelCoordToChunkCoordInvalidChunkSize(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	chunk.VoxelCoordToChunkCoord(chunk.VoxelCoordinate{0, 0, 0}, 0)
}

func TestVoxelCoordToChunkCoord(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		voxelCoord       chunk.VoxelCoordinate
		expectChunkCoord chunk.Position
		chunkSize        uint32
	}{
		{
			voxelCoord:       chunk.VoxelCoordinate{0, 0, 0},
			expectChunkCoord: chunk.Position{0, 0, 0},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, 0, 0},
			expectChunkCoord: chunk.Position{1, 0, 0},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, -1, 3},
			expectChunkCoord: chunk.Position{1, -1, 3},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, 1, 1},
			expectChunkCoord: chunk.Position{0, 0, 0},
			chunkSize:        2,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{-3, -1, -2},
			expectChunkCoord: chunk.Position{-1, -1, -1},
			chunkSize:        3,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{29, 29, 29},
			expectChunkCoord: chunk.Position{2, 2, 2},
			chunkSize:        10,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{30, 30, 30},
			expectChunkCoord: chunk.Position{3, 3, 3},
			chunkSize:        10,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(fmt.Sprintf("%v should be in chunk %v with size %v\n",
			tC.voxelCoord, tC.expectChunkCoord, tC.chunkSize), func(t *testing.T) {
			t.Parallel()
			actualChunkCoord := chunk.VoxelCoordToChunkCoord(tC.voxelCoord, tC.chunkSize)
			if actualChunkCoord != tC.expectChunkCoord {
				t.Fatalf("expected chunk coord to be %v but was %v", tC.expectChunkCoord, actualChunkCoord)
			}
		})
	}
}

func TestLoadChunkFromFlatDataDetectsError(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		badData   []float32
		chunkSize uint32
		chPos     chunk.Position
	}{
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 2},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 2, 3},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 2, 3, 4},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 2, 3, 4, 5},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 2, 3, 4, 5, 6},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits + 1), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll + 1)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{1, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{0, -1, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{0, 0, 0},
			badData:   []float32{0, 0, 2, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.Position{1, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos: chunk.Position{0, 0, 0},
			badData: []float32{ // order swap
				0, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll),
				1, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll),
				0, 0, 1, float32(chunk.LargestVbits), float32(chunk.LightAll),
				1, 0, 1, float32(chunk.LargestVbits), float32(chunk.LightAll),
				0, 1, 0, float32(chunk.LargestVbits), float32(chunk.LightAll),
				1, 1, 0, float32(chunk.LargestVbits), float32(chunk.LightAll),
				0, 1, 1, float32(chunk.LargestVbits), float32(chunk.LightAll),
				1, 1, 1, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 2,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run("", func(t *testing.T) {
			t.Parallel()
			defer func() {
				if err := recover(); err == nil {
					t.Fatalf("did not panic for bad data: %v", tC.badData)
				}
			}()
			chunk.NewFromData(tC.badData, tC.chunkSize, tC.chPos)
		})
	}
}

func TestValidChunkDataScenario(t *testing.T) {
	t.Parallel()
	ch := chunk.NewEmpty(chunk.Position{1, 1, 1}, 3)
	ch.SetBlockType(chunk.VoxelCoordinate{3, 3, 3}, chunk.BlockTypeDirt)
	ch.SetAdjacency(chunk.VoxelCoordinate{3, 4, 3}, chunk.AdjacentBack)
	ch.SetLighting(chunk.VoxelCoordinate{4, 4, 4}, chunk.LightTop, 6)
	data := ch.GetFlatData()
	chFromData := chunk.NewFromData(data, ch.Size(), chunk.Position{1, 1, 1})
	actualBlockType := chFromData.BlockType(chunk.VoxelCoordinate{3, 3, 3})
	actualAdjacency := chFromData.Adjacency(chunk.VoxelCoordinate{3, 4, 3})
	actualLighting := chFromData.Lighting(chunk.VoxelCoordinate{4, 4, 4}, chunk.LightTop)
	if actualBlockType != chunk.BlockTypeDirt {
		t.Fatal("recovered wrong block type")
	}
	if actualAdjacency != chunk.AdjacentBack {
		t.Fatal("recovered wrong adjacency")
	}
	if actualLighting != 6 {
		t.Fatal("recovered wrong lighting")
	}
}
