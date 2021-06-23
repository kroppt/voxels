package world

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// metadata file:
// numChunks|x|y|z|offset|x|y|z|offset|x|y|z|offset
// chunk data file:
// vx|vy|vz|vbits|vx|vy|vz|vbits|vx|vy|vz|vbits|
type Cache struct {
	chunksToBeWritten map[ChunkPos]*Chunk
	metaFile          *os.File
	dataFile          *os.File
	numChunks         int32
	writeThreshold    int32
}

func NewCache(metaFileName, dataFileName string, writeThreshold int32) (*Cache, error) {
	meta, err := os.OpenFile(metaFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	data, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &Cache{
		metaFile:          meta,
		dataFile:          data,
		chunksToBeWritten: make(map[ChunkPos]*Chunk),
		writeThreshold:    writeThreshold,
	}, nil
}

func (c *Cache) Save(ch *Chunk) {
	c.chunksToBeWritten[ch.Pos] = ch
	c.numChunks++
	if c.numChunks >= c.writeThreshold {
		c.writeBufferToFile()
	}
}

func (c *Cache) Load(pos ChunkPos) (*Chunk, bool) {
	ch, ok := c.chunksToBeWritten[pos]
	if ok {
		ch.dirty = true
		return ch, true
	} else {
		dataOff, found := c.findChunkVoxelDataOffset(pos)
		if found {
			flatData := c.readChunkVoxelData(int64(dataOff))
			return NewChunkLoaded(ChunkSize, pos, flatData), true
		}
	}
	return nil, false
}

func readChunkMetadata(f *os.File, byteOff int64) (int, int, int, int32, error) {
	var x, y, z, offset int32
	b := make([]byte, 16)
	n, err := f.ReadAt(b, byteOff)
	if err != nil || n != 16 {
		return 0, 0, 0, 0, err
	}
	readBuf := bytes.NewBuffer(b)
	vals := []interface{}{&x, &y, &z, &offset}
	for _, val := range vals {
		err = binary.Read(readBuf, binary.BigEndian, val)
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}
	return int(x), int(y), int(z), offset, nil
}

func (c *Cache) readChunkVoxelData(byteOff int64) []int32 {
	size := ChunkSize * ChunkSize * ChunkSize * 16
	b := make([]byte, size)
	n, err := c.dataFile.ReadAt(b, byteOff)
	if err != nil || n != size {
		panic(err)
	}

	readBuf := bytes.NewBuffer(b)
	numInts := ChunkSize * ChunkSize * ChunkSize * 4
	flatData := make([]int32, numInts)
	for i := 0; i < numInts; i++ {
		err = binary.Read(readBuf, binary.BigEndian, &flatData[i])
		if err != nil {
			panic(err)
		}
	}
	return flatData
}

func (c *Cache) writeChunkMetadata(byteOff int64, p ChunkPos, dataOff int32) error {
	var writeBuf bytes.Buffer
	vals := []interface{}{int32(p.X), int32(p.Y), int32(p.Z), int32(dataOff)}
	for _, val := range vals {
		writeErr := binary.Write(&writeBuf, binary.BigEndian, val)
		if writeErr != nil {
			return writeErr
		}
	}
	n, err := c.metaFile.WriteAt(writeBuf.Bytes(), byteOff)
	if err != nil || n != 16 {
		return fmt.Errorf("failed to write metadata at byteOff %v", byteOff)
	}
	return nil
}

func (c *Cache) writeChunkVoxelData(byteOff int64, ch *Chunk) error {
	var writeBuf bytes.Buffer
	maxIdx := 4 * ChunkSize * ChunkSize * ChunkSize
	for i := 0; i < maxIdx; i++ {
		writeErr := binary.Write(&writeBuf, binary.BigEndian, int32(ch.flatData[i]))
		if writeErr != nil {
			return writeErr
		}
	}
	n, err := c.dataFile.WriteAt(writeBuf.Bytes(), int64(byteOff))
	if err != nil || n != 4*maxIdx {
		return err
	}
	return nil
}

type ChunkOffs struct {
	metaOff int32
	dataOff int32
}

func (c *Cache) writeBufferToFile() {
	if c.numChunks == 0 {
		return
	}
	numFileChunks, exists := c.GetNumChunksInFile()
	if !exists {
		// meta file did not exist
		var writeBuf bytes.Buffer
		writeErr := binary.Write(&writeBuf, binary.BigEndian, c.numChunks)
		if writeErr != nil {
			panic(writeErr)
		}
		c.metaFile.WriteAt(writeBuf.Bytes(), 0)
	}

	fileChunks := make(map[ChunkPos]ChunkOffs)
	for i := 0; i < int(numFileChunks); i++ {
		metaOff := int64(4 + 16*i)
		x, y, z, off, err := readChunkMetadata(c.metaFile, metaOff)
		if err == io.EOF {
			panic("metadata file had fewer chunk meta datas than reported by the file")
		}
		fileChunks[ChunkPos{x, y, z}] = ChunkOffs{
			metaOff: int32(metaOff),
			dataOff: off,
		}
	}

	metaEndIdx := int64(4 + numFileChunks*16)
	chunkDataSize := int64(ChunkSize * ChunkSize * ChunkSize * 16)
	dataEndIdx := int64(numFileChunks) * chunkDataSize
	for pos, chunk := range c.chunksToBeWritten {
		if off, ok := fileChunks[pos]; ok {
			// chunk existed in file, overwrite
			// overwrite chunk's meta data at offset
			err := c.writeChunkMetadata(int64(off.metaOff), pos, off.dataOff)
			if err != nil {
				panic("error overwriting existing chunk metadata")
			}
			// overwrite chunk's voxel data at offset
			err = c.writeChunkVoxelData(int64(off.dataOff), chunk)
			if err != nil {
				panic("error overwriting existing chunk voxel data")
			}
		} else {
			// new chunk
			numFileChunks++
			// append chunk metadata
			err := c.writeChunkMetadata(int64(metaEndIdx), pos, int32(dataEndIdx))
			if err != nil {
				panic("error appending new chunk metadata")
			}
			metaEndIdx += 16
			// append chunk voxel data
			err = c.writeChunkVoxelData(int64(dataEndIdx), chunk)
			if err != nil {
				panic("error appending new chunk voxel data")
			}
			dataEndIdx += chunkDataSize
		}

	}

	var writeBuf bytes.Buffer
	writeErr := binary.Write(&writeBuf, binary.BigEndian, int32(numFileChunks))
	if writeErr != nil {
		panic(writeErr)
	}
	c.metaFile.WriteAt(writeBuf.Bytes(), 0)

	c.numChunks = 0
	c.chunksToBeWritten = make(map[ChunkPos]*Chunk)
}

func (c *Cache) GetNumChunksInFile() (int32, bool) {
	var numChunks int32
	b := make([]byte, 4)
	n, err := c.metaFile.ReadAt(b, 0)
	if err != nil || n != 4 {
		return 0, false
	}
	// file already existed
	readBuf := bytes.NewBuffer(b)
	err = binary.Read(readBuf, binary.BigEndian, &numChunks)
	if err != nil {
		panic(err)
	}
	return numChunks, true
}

func (c *Cache) findChunkVoxelDataOffset(pos ChunkPos) (int32, bool) {
	numChunks, received := c.GetNumChunksInFile()
	if !received {
		return 0, false
	}

	for i := 0; i < int(numChunks); i++ {
		metaOff := int64(4 + 16*i)
		x, y, z, off, err := readChunkMetadata(c.metaFile, metaOff)
		if err == io.EOF {
			// fmt.Printf("%v %v %v %v", x, y, z, off)
			panic(err)
		}
		if x == pos.X && y == pos.Y && z == pos.Z {
			return off, true
		}
	}
	return 0, false
}

func (c *Cache) Sync() {
	c.writeBufferToFile()
}

func (c *Cache) Destroy() {
	c.metaFile.Close()
	c.dataFile.Close()
}
