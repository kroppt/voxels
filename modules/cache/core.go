package cache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/spf13/afero"
)

type core struct {
	voxelFile    afero.File
	chunkFile    afero.File
	regionFile   afero.File
	settingsRepo settings.Interface
}

type regionPosition struct {
	x int32
	y int32
	z int32
}

func chunkPosToRegionPos(pos chunk.Position, regionSize uint32) regionPosition {
	if regionSize == 0 {
		panic("regionSize 0 is invalid")
	}
	x, y, z := pos.X, pos.Y, pos.Z
	size := int32(regionSize)
	if pos.X < 0 {
		x++
	}
	if pos.Y < 0 {
		y++
	}
	if pos.Z < 0 {
		z++
	}
	x /= size
	y /= size
	z /= size
	if pos.X < 0 {
		x--
	}
	if pos.Y < 0 {
		y--
	}
	if pos.Z < 0 {
		z--
	}
	return regionPosition{x: x, y: y, z: z}

}

func (c *core) save(ch chunk.Chunk) {
	regionPos := chunkPosToRegionPos(ch.Position(), c.settingsRepo.GetRegionSize())
	regionIdx, ok := c.getRegionIdx(regionPos)
	if !ok {
		// region wasn't registered
		info, err := c.regionFile.Stat()
		if err != nil {
			log.Print(err)
		}
		regionOff := info.Size()
		info, err = c.chunkFile.Stat()
		if err != nil {
			log.Print(err)
		}
		regionIdx := info.Size()
		c.writeRegionFileAt(regionPos.x, regionPos.y, regionPos.z, int32(regionIdx), regionOff)
		c.writeRegionAt(regionPos, regionIdx)
		// chunk couldn't have been registered because the region wasn't
		info, err = c.voxelFile.Stat()
		if err != nil {
			log.Print(err)
		}
		chunkIdx := int32(info.Size())
		c.writeChunkFileAt(chunkIdx, int64(int32(regionIdx)+4*chunkPosToDataOffset(ch.Position(), regionPos, int32(c.settingsRepo.GetRegionSize()))))
		c.writeChunkAt(ch, int64(chunkIdx))

	} else {
		// region was registered, but maybe chunk wasnt
		chunkIdx := c.getChunkIdx(regionIdx + 4*chunkPosToDataOffset(ch.Position(), regionPos, int32(c.settingsRepo.GetRegionSize())))
		if chunkIdx == -1 {
			// chunk was not registered
			info, err := c.voxelFile.Stat()
			if err != nil {
				log.Print(err)
			}
			chunkIdx = int32(info.Size())
			c.writeChunkFileAt(chunkIdx,
				int64(regionIdx+4*chunkPosToDataOffset(ch.Position(), regionPos, int32(c.settingsRepo.GetRegionSize()))))
			c.writeChunkAt(ch, int64(chunkIdx))
		} else {
			// both region and chunk were registered
			c.writeChunkAt(ch, int64(chunkIdx))
		}
	}
}

func (c *core) load(pos chunk.Position) (chunk.Chunk, bool) {
	regionPos := chunkPosToRegionPos(pos, c.settingsRepo.GetRegionSize())
	regionIdx, ok := c.getRegionIdx(regionPos)
	if !ok {
		return chunk.Chunk{}, false
	} else {
		// region existed, chunk registered?
		chunkIdx := c.getChunkIdx(regionIdx + 4*chunkPosToDataOffset(pos, regionPos, int32(c.settingsRepo.GetRegionSize())))
		if chunkIdx == -1 {
			return chunk.Chunk{}, false
		}

		chunkSize := c.settingsRepo.GetChunkSize()
		numElems := chunk.VertSize * chunkSize * chunkSize * chunkSize
		byteSize := chunk.BytesPerElement * numElems
		bs := make([]byte, byteSize)

		n, err := c.voxelFile.ReadAt(bs, int64(chunkIdx))
		if uint32(n) != byteSize {
			log.Printf("(load) expected %v bytes to be read, but only %v were read", byteSize, n)
			return chunk.Chunk{}, false
		}
		if !errors.Is(err, io.EOF) && err != nil {
			log.Print(err)
			return chunk.Chunk{}, false
		}
		buf := bytes.NewBuffer(bs)
		flatData := make([]float32, numElems)
		err = binary.Read(buf, binary.LittleEndian, flatData)
		if err != nil {
			log.Print(err)
			return chunk.Chunk{}, false
		}
		return chunk.NewFromData(flatData, chunkSize, pos), true
	}

}

func (c *core) writeChunkAt(ch chunk.Chunk, off int64) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, ch.GetFlatData())
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.voxelFile.WriteAt(buf.Bytes(), off)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("(write) expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func chunkPosToDataOffset(chunkPos chunk.Position, regionPos regionPosition, size int32) int32 {
	i := chunkPos.X - regionPos.x*size
	j := chunkPos.Y - regionPos.y*size
	k := chunkPos.Z - regionPos.z*size
	return i + j*size + k*size*size
}

func (c *core) getEmptyRegionData(regionPos regionPosition) []int32 {
	size := int32(c.settingsRepo.GetRegionSize())
	data := make([]int32, size*size*size)
	for x := regionPos.x * size; x < regionPos.x*size+size; x++ {
		for y := regionPos.y * size; y < regionPos.y*size+size; y++ {
			for z := regionPos.z * size; z < regionPos.z*size+size; z++ {
				off := chunkPosToDataOffset(chunk.Position{X: x, Y: y, Z: z}, regionPos, size)
				data[off] = -1.0
			}
		}
	}
	return data
}

func (c *core) writeRegionAt(regionPos regionPosition, off int64) {
	var buf bytes.Buffer
	emptyData := c.getEmptyRegionData(regionPos)
	err := binary.Write(&buf, binary.LittleEndian, emptyData)
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.chunkFile.WriteAt(buf.Bytes(), off)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("(writeRegionAt) expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (c *core) writeRegionFileAt(x, y, z, metaIdx int32, off int64) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, []int32{x, y, z, metaIdx})
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.regionFile.WriteAt(buf.Bytes(), off)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("(writeFileAt) expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (c *core) writeChunkFileAt(metaIdx int32, off int64) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, metaIdx)
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.chunkFile.WriteAt(buf.Bytes(), off)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("(writeFileAt) expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (c *core) getChunkIdx(off int32) int32 {
	byteSize := 4
	bs := make([]byte, byteSize)
	n, err := c.chunkFile.ReadAt(bs, int64(off))
	if n != byteSize {
		log.Printf("(getChunkIdx) expected %v bytes to be read, but only %v were read", byteSize, n)
		return -1
	}
	if !errors.Is(err, io.EOF) && err != nil {
		log.Print(err)
		return -1
	}
	buf := bytes.NewBuffer(bs)
	var chunkIdx int32
	err = binary.Read(buf, binary.LittleEndian, &chunkIdx)
	if err != nil {
		log.Print(err)
		return -1
	}
	return chunkIdx
}

func (c *core) getRegionIdx(regionPos regionPosition) (int32, bool) {
	curr := int64(0)
	for {
		numElems := 4
		byteSize := 4 * numElems
		bs := make([]byte, byteSize)
		n, err := c.regionFile.ReadAt(bs, curr)
		if n == 0 && errors.Is(err, io.EOF) {
			return -1, false
		}
		if n != byteSize {
			log.Printf("(getRegionIdx) expected %v bytes to be read, but only %v were read", byteSize, n)
			return -1, false
		}
		if !errors.Is(err, io.EOF) && err != nil {
			log.Print(err)
			return -1, false
		}
		buf := bytes.NewBuffer(bs)
		regionEntry := make([]int32, numElems)
		err = binary.Read(buf, binary.LittleEndian, regionEntry)
		if err != nil {
			log.Print(err)
			return -1, false
		}
		key := regionPosition{
			x: regionEntry[0],
			y: regionEntry[1],
			z: regionEntry[2],
		}
		if key == regionPos {
			return regionEntry[3], true
		}
		curr += int64(byteSize)
	}
}

func (c *core) close() {
	err := c.voxelFile.Close()
	if err != nil {
		panic(err)
	}
	err = c.chunkFile.Close()
	if err != nil {
		panic(err)
	}
	err = c.regionFile.Close()
	if err != nil {
		panic(err)
	}
}
