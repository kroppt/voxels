package chunk_test

import (
	"container/list"
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
		chunk.NewChunkEmpty(chunk.ChunkCoordinate{}, 0)
	})
	t.Run("cannot set block type out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetBlockType(chunk.VoxelCoordinate{2, 2, 2}, chunk.BlockTypeAir)
	})
	t.Run("cannot set block type out of chunk bounds complex coords", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{-1, -1, -1}, 2)
		c.SetBlockType(chunk.VoxelCoordinate{-1, -1, 1}, chunk.BlockTypeAir)
	})
	t.Run("cannot get block type out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.BlockType(chunk.VoxelCoordinate{2, 2, 2})
	})
	t.Run("cannot set lighting out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{2, 2, 2}, chunk.LightLeft, 5)
	})
	t.Run("cannot get lighting out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.Lighting(chunk.VoxelCoordinate{2, 2, 2}, chunk.LightFront)
	})
	t.Run("cannot set adjacency out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetAdjacency(chunk.VoxelCoordinate{2, 2, 2}, chunk.AdjacentBack)
	})
	t.Run("cannot get adjacency out of chunk bounds", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.Adjacency(chunk.VoxelCoordinate{2, 2, 2})
	})
	t.Run("cannot set invalid adjacency", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetAdjacency(chunk.VoxelCoordinate{0, 0, 0}, 0b1000000)
	})
	t.Run("cannot set invalid light intensity", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{0, 0, 0}, chunk.LightBack, 16)
	})
	t.Run("cannot specifcy invalid light face when setting", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.SetLighting(chunk.VoxelCoordinate{0, 0, 0}, 21, 5)
	})
	t.Run("cannot specifcy invalid light face when getting", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic, but didn't")
			}
		}()
		c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 2)
		c.Lighting(chunk.VoxelCoordinate{0, 0, 0}, 21)
	})
}

func TestChunk(t *testing.T) {
	t.Parallel()
	t.Run("test get chunk position", func(t *testing.T) {
		t.Parallel()
		expected := chunk.ChunkCoordinate{1, 2, 3}
		chunk := chunk.NewChunkEmpty(expected, 10)
		actual := chunk.Position()
		if actual != expected {
			t.Fatalf("expected to get chunk position %v but got %v", expected, actual)
		}
	})
	t.Run("test get chunk size", func(t *testing.T) {
		t.Parallel()
		expected := uint32(10)
		chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{1, 2, 3}, expected)
		actual := chunk.Size()
		if actual != expected {
			t.Fatalf("expected to get chunk size %v but got %v", expected, actual)
		}
	})
	t.Run("check that default block type is air", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeAir
		voxelCoordinate := chunk.VoxelCoordinate{4, 4, 4}
		chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 6)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("set block type of one voxel to dirt and get it back", func(t *testing.T) {
		t.Parallel()
		expected := chunk.BlockTypeDirt
		voxelCoordinate := chunk.VoxelCoordinate{5, 5, 5}
		chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 10)
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
		chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{1, 1, 1}, 10)
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
		chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{-1, -1, -1}, 10)
		chunk.SetBlockType(voxelCoordinate, expected)
		actual := chunk.BlockType(voxelCoordinate)
		if actual != expected {
			t.Fatalf("expected to get back block type of %v but got %v", expected, actual)
		}
	})
	t.Run("flat data contains correct indices of voxels", func(t *testing.T) {
		t.Parallel()
		chPos := chunk.ChunkCoordinate{0, 0, 0}
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
		ch := chunk.NewChunkEmpty(chPos, uint32(size))
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
			chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 10)
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
			chunk := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 10)
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
		expectChunkCoord chunk.ChunkCoordinate
		chunkSize        uint32
	}{
		{
			voxelCoord:       chunk.VoxelCoordinate{0, 0, 0},
			expectChunkCoord: chunk.ChunkCoordinate{0, 0, 0},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, 0, 0},
			expectChunkCoord: chunk.ChunkCoordinate{1, 0, 0},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, -1, 3},
			expectChunkCoord: chunk.ChunkCoordinate{1, -1, 3},
			chunkSize:        1,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{1, 1, 1},
			expectChunkCoord: chunk.ChunkCoordinate{0, 0, 0},
			chunkSize:        2,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{-3, -1, -2},
			expectChunkCoord: chunk.ChunkCoordinate{-1, -1, -1},
			chunkSize:        3,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{29, 29, 29},
			expectChunkCoord: chunk.ChunkCoordinate{2, 2, 2},
			chunkSize:        10,
		},
		{
			voxelCoord:       chunk.VoxelCoordinate{30, 30, 30},
			expectChunkCoord: chunk.ChunkCoordinate{3, 3, 3},
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
		chPos     chunk.ChunkCoordinate
	}{
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 2},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 2, 3},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 2, 3, 4},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 2, 3, 4, 5},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 2, 3, 4, 5, 6},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits + 1), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll + 1)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{1, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{0, -1, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{0, 0, 0},
			badData:   []float32{0, 0, 2, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos:     chunk.ChunkCoordinate{1, 0, 0},
			badData:   []float32{0, 0, 0, float32(chunk.LargestVbits), float32(chunk.LightAll)},
			chunkSize: 1,
		},
		{
			chPos: chunk.ChunkCoordinate{0, 0, 0},
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
			chunk.NewChunkFromData(tC.badData, tC.chunkSize, tC.chPos)
		})
	}
}

func TestValidChunkDataScenario(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{1, 1, 1}, 3)
	ch.SetBlockType(chunk.VoxelCoordinate{3, 3, 3}, chunk.BlockTypeDirt)
	ch.SetAdjacency(chunk.VoxelCoordinate{3, 4, 3}, chunk.AdjacentBack)
	ch.SetLighting(chunk.VoxelCoordinate{4, 4, 4}, chunk.LightTop, 6)
	data := ch.GetFlatData()
	chFromData := chunk.NewChunkFromData(data, ch.Size(), chunk.ChunkCoordinate{1, 1, 1})
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

func TestForEachVoxelInChunk(t *testing.T) {
	t.Parallel()
	chPos := chunk.ChunkCoordinate{-1, 2, 3}
	chSize := 2
	expected := map[chunk.VoxelCoordinate]struct{}{
		{-2, 4, 6}: {},
		{-1, 4, 6}: {},
		{-2, 5, 6}: {},
		{-1, 5, 6}: {},
		{-2, 4, 7}: {},
		{-1, 4, 7}: {},
		{-2, 5, 7}: {},
		{-1, 5, 7}: {},
	}
	ch := chunk.NewChunkEmpty(chPos, uint32(chSize))
	ch.ForEachVoxel(func(voxPos chunk.VoxelCoordinate) {
		if _, ok := expected[voxPos]; !ok {
			t.Fatalf("expected %v to be in map, but wasn't", voxPos)
		}
	})
}

func TestGetVbitsForVoxel(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{1, 1, 1}, 3)
	ch.SetBlockType(chunk.VoxelCoordinate{3, 3, 3}, chunk.BlockTypeCorrupted)
	ch.SetAdjacency(chunk.VoxelCoordinate{3, 3, 3}, chunk.AdjacentBack|chunk.AdjacentBottom|chunk.AdjacentY)
	expectedVbits := uint32(chunk.BlockTypeCorrupted<<6) | uint32(chunk.AdjacentBack|chunk.AdjacentBottom|chunk.AdjacentY)
	actualVbits := ch.Vbits(chunk.VoxelCoordinate{3, 3, 3})
	if actualVbits != expectedVbits {
		t.Fatalf("expected vbits %v but got %v", expectedVbits, actualVbits)
	}
}

func TestAddAdjacency(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront|chunk.AdjacentBack|chunk.AdjacentRight)
	ch.AddAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentLeft)
	expected := chunk.AdjacentFront | chunk.AdjacentBack | chunk.AdjacentRight | chunk.AdjacentLeft
	actual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 1})
	if actual != expected {
		t.Fatalf("expected to get adjacency %v but got %v", expected, actual)
	}
}

func TestAddSameAdjacencyNoChange(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront)
	ch.AddAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront)
	expected := chunk.AdjacentFront
	actual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 1})
	if actual != expected {
		t.Fatalf("expected to get adjacency %v but got %v", expected, actual)
	}
}

func TestAddAdjacencyFromNone(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	ch.AddAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentBottom)
	expected := chunk.AdjacentBottom
	actual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 1})
	if actual != expected {
		t.Fatalf("expected to get adjacency %v but got %v", expected, actual)
	}
}

func TestRemoveAdjacency(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront|chunk.AdjacentBack|chunk.AdjacentRight)
	ch.RemoveAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront)
	expected := chunk.AdjacentBack | chunk.AdjacentRight
	actual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 1})
	if actual != expected {
		t.Fatalf("expected to get adjacency %v but got %v", expected, actual)
	}
}

func TestRemoveAdjacencyFromNoneNoChange(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	ch.RemoveAdjacency(chunk.VoxelCoordinate{1, 1, 1}, chunk.AdjacentFront)
	expected := chunk.AdjacentMask(chunk.AdjacentNone)
	actual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 1})
	if actual != expected {
		t.Fatalf("expected to get adjacency %v but got %v", expected, actual)
	}
}

func TestSetBlockTypeAddsAdjacenciesAutomatically(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	additionalChanges := ch.SetBlockType(chunk.VoxelCoordinate{1, 1, 1}, chunk.BlockTypeDirt)
	if additionalChanges.Len() != 0 {
		t.Fatal("received pending actions when there should be none")
	}
	frontExpect := chunk.AdjacentBack
	frontActual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 0})
	backExpect := chunk.AdjacentFront
	backActual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 2})
	rightExpect := chunk.AdjacentLeft
	rightActual := ch.Adjacency(chunk.VoxelCoordinate{2, 1, 1})
	leftExpect := chunk.AdjacentRight
	leftActual := ch.Adjacency(chunk.VoxelCoordinate{0, 1, 1})
	upExpect := chunk.AdjacentBottom
	upActual := ch.Adjacency(chunk.VoxelCoordinate{1, 2, 1})
	downExpect := chunk.AdjacentTop
	downActual := ch.Adjacency(chunk.VoxelCoordinate{1, 0, 1})
	if frontActual != frontExpect {
		t.Fatalf("expected front adjacency %v but got %v", frontExpect, frontActual)
	}
	if backActual != backExpect {
		t.Fatalf("expected back adjacency %v but got %v", backExpect, backActual)
	}
	if rightActual != rightExpect {
		t.Fatalf("expected right adjacency %v but got %v", rightExpect, rightActual)
	}
	if leftActual != leftExpect {
		t.Fatalf("expected left adjacency %v but got %v", leftExpect, leftActual)
	}
	if upActual != upExpect {
		t.Fatalf("expected up adjacency %v but got %v", upExpect, upActual)
	}
	if downActual != downExpect {
		t.Fatalf("expected down adjacency %v but got %v", downExpect, downActual)
	}
}

func TestSetBlockTypeRemovesAdjacenciesAutomaticallyWithinChunk(t *testing.T) {
	t.Parallel()
	ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
	expect := chunk.AdjacentMask(chunk.AdjacentNone)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 1, 0}, chunk.AdjacentBack)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 1, 2}, chunk.AdjacentFront)
	ch.SetAdjacency(chunk.VoxelCoordinate{2, 1, 1}, chunk.AdjacentLeft)
	ch.SetAdjacency(chunk.VoxelCoordinate{0, 1, 1}, chunk.AdjacentRight)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 2, 1}, chunk.AdjacentBottom)
	ch.SetAdjacency(chunk.VoxelCoordinate{1, 0, 1}, chunk.AdjacentTop)

	additionalChanges := ch.SetBlockType(chunk.VoxelCoordinate{1, 1, 1}, chunk.BlockTypeAir)
	if additionalChanges.Len() != 0 {
		t.Fatal("received pending actions when there should be none")
	}

	frontActual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 0})
	backActual := ch.Adjacency(chunk.VoxelCoordinate{1, 1, 2})
	rightActual := ch.Adjacency(chunk.VoxelCoordinate{2, 1, 1})
	leftActual := ch.Adjacency(chunk.VoxelCoordinate{0, 1, 1})
	upActual := ch.Adjacency(chunk.VoxelCoordinate{1, 2, 1})
	downActual := ch.Adjacency(chunk.VoxelCoordinate{1, 0, 1})

	if frontActual != expect {
		t.Fatalf("expected front adjacency %v but got %v", expect, frontActual)
	}
	if backActual != expect {
		t.Fatalf("expected back adjacency %v but got %v", expect, backActual)
	}
	if rightActual != expect {
		t.Fatalf("expected right adjacency %v but got %v", expect, rightActual)
	}
	if leftActual != expect {
		t.Fatalf("expected left adjacency %v but got %v", expect, leftActual)
	}
	if upActual != expect {
		t.Fatalf("expected up adjacency %v but got %v", expect, upActual)
	}
	if downActual != expect {
		t.Fatalf("expected down adjacency %v but got %v", expect, downActual)
	}
}

func TestRemoveBlockAtChunkBoundaries(t *testing.T) {
	t.Parallel()
	t.Run("simple removal one chunk affected", func(t *testing.T) {
		ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 3)
		expected := list.New()
		expected.PushBack(chunk.PendingAction{
			ChPos:    chunk.ChunkCoordinate{0, 0, -1},
			VoxPos:   chunk.VoxelCoordinate{1, 1, -1},
			HideFace: false,
			Face:     chunk.AdjacentBack,
		})
		actual := ch.SetBlockType(chunk.VoxelCoordinate{1, 1, 0}, chunk.BlockTypeAir)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("received wrong pending actions")
		}
	})
	t.Run("extreme removal edge case", func(t *testing.T) {
		ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 1)
		expected := list.New()
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 0, -1},
			VoxPos: chunk.VoxelCoordinate{0, 0, -1}, HideFace: false, Face: chunk.AdjacentBack})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 0, 1},
			VoxPos: chunk.VoxelCoordinate{0, 0, 1}, HideFace: false, Face: chunk.AdjacentFront})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{-1, 0, 0},
			VoxPos: chunk.VoxelCoordinate{-1, 0, 0}, HideFace: false, Face: chunk.AdjacentRight})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{1, 0, 0},
			VoxPos: chunk.VoxelCoordinate{1, 0, 0}, HideFace: false, Face: chunk.AdjacentLeft})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, -1, 0},
			VoxPos: chunk.VoxelCoordinate{0, -1, 0}, HideFace: false, Face: chunk.AdjacentTop})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 1, 0},
			VoxPos: chunk.VoxelCoordinate{0, 1, 0}, HideFace: false, Face: chunk.AdjacentBottom})
		actual := ch.SetBlockType(chunk.VoxelCoordinate{0, 0, 0}, chunk.BlockTypeAir)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("received wrong pending actions")
		}
	})
	t.Run("extreme add edge case", func(t *testing.T) {
		ch := chunk.NewChunkEmpty(chunk.ChunkCoordinate{0, 0, 0}, 1)
		expected := list.New()
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 0, -1},
			VoxPos: chunk.VoxelCoordinate{0, 0, -1}, HideFace: true, Face: chunk.AdjacentBack})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 0, 1},
			VoxPos: chunk.VoxelCoordinate{0, 0, 1}, HideFace: true, Face: chunk.AdjacentFront})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{-1, 0, 0},
			VoxPos: chunk.VoxelCoordinate{-1, 0, 0}, HideFace: true, Face: chunk.AdjacentRight})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{1, 0, 0},
			VoxPos: chunk.VoxelCoordinate{1, 0, 0}, HideFace: true, Face: chunk.AdjacentLeft})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, -1, 0},
			VoxPos: chunk.VoxelCoordinate{0, -1, 0}, HideFace: true, Face: chunk.AdjacentTop})
		expected.PushBack(chunk.PendingAction{ChPos: chunk.ChunkCoordinate{0, 1, 0},
			VoxPos: chunk.VoxelCoordinate{0, 1, 0}, HideFace: true, Face: chunk.AdjacentBottom})
		actual := ch.SetBlockType(chunk.VoxelCoordinate{0, 0, 0}, chunk.BlockTypeDirt)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("received wrong pending actions")
		}
	})
}

func TestChunkAppliesActions(t *testing.T) {
	c := chunk.NewChunkEmpty(chunk.ChunkCoordinate{}, 2)
	actions := list.New()
	expect1 := chunk.AdjacentAll
	expect2 := chunk.AdjacentX
	actions.PushBack(chunk.PendingAction{
		ChPos:    chunk.ChunkCoordinate{},
		VoxPos:   chunk.VoxelCoordinate{0, 0, 0},
		HideFace: true,
		Face:     chunk.AdjacentAll,
	})
	actions.PushBack(chunk.PendingAction{
		ChPos:    chunk.ChunkCoordinate{},
		VoxPos:   chunk.VoxelCoordinate{1, 1, 1},
		HideFace: true,
		Face:     chunk.AdjacentX,
	})
	c.ApplyActions(actions)
	actual1 := c.Adjacency(chunk.VoxelCoordinate{0, 0, 0})
	actual2 := c.Adjacency(chunk.VoxelCoordinate{1, 1, 1})

	if actual1 != expect1 {
		t.Fatalf("(1) expected adjacency %v but got %v", expect1, actual1)
	}
	if actual2 != expect2 {
		t.Fatalf("(2) expected adjacency %v but got %v", expect2, actual2)
	}
}
