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
	dataFile     afero.File
	metaFile     afero.File
	settingsRepo settings.Interface
}

func (c *core) save(ch chunk.Chunk) {
	metaIdx, ok := c.getMetaIdx(ch.Position())
	if !ok {
		info, err := c.metaFile.Stat()
		if err != nil {
			log.Print(err)
		}
		metaOff := info.Size()
		info, err = c.dataFile.Stat()
		if err != nil {
			log.Print(err)
		}
		metaIdx := info.Size()
		c.writeMetaDataAt(ch.Position(), int32(metaIdx), metaOff)
		c.writeChunkAt(ch, metaIdx)
		return
	}
	c.writeChunkAt(ch, int64(metaIdx))
}

func (c *core) load(pos chunk.Position) (chunk.Chunk, bool) {
	metaIdx, ok := c.getMetaIdx(pos)
	if !ok {
		return chunk.Chunk{}, false
	}

	chunkSize := c.settingsRepo.GetChunkSize()
	numElems := chunk.VertSize * chunkSize * chunkSize * chunkSize
	byteSize := chunk.BytesPerElement * numElems
	bs := make([]byte, byteSize)
	n, err := c.dataFile.ReadAt(bs, int64(metaIdx))
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

func (c *core) writeChunkAt(ch chunk.Chunk, off int64) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, ch.GetFlatData())
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.dataFile.WriteAt(buf.Bytes(), off)
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

func (c *core) writeMetaDataAt(pos chunk.Position, metaIdx int32, off int64) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, []int32{pos.X, pos.Y, pos.Z, metaIdx})
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.metaFile.WriteAt(buf.Bytes(), off)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("(write meta) expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (c *core) getMetaIdx(pos chunk.Position) (int32, bool) {
	curr := int64(0)
	for {
		numElems := 4
		byteSize := 4 * numElems
		bs := make([]byte, byteSize)
		n, err := c.metaFile.ReadAt(bs, curr)
		if n == 0 && errors.Is(err, io.EOF) {
			return -1, false
		}
		if n != byteSize {
			log.Printf("(get meta) expected %v bytes to be read, but only %v were read", byteSize, n)
			return -1, false
		}
		if !errors.Is(err, io.EOF) && err != nil {
			log.Print(err)
			return -1, false
		}
		buf := bytes.NewBuffer(bs)
		metaDataEntry := make([]int32, numElems)
		err = binary.Read(buf, binary.LittleEndian, metaDataEntry)
		if err != nil {
			log.Print(err)
			return -1, false
		}
		key := chunk.Position{
			X: metaDataEntry[0],
			Y: metaDataEntry[1],
			Z: metaDataEntry[2],
		}
		if key == pos {
			return metaDataEntry[3], true
		}
		curr += int64(byteSize)
	}
}

func (c *core) close() {
	err := c.dataFile.Close()
	if err != nil {
		panic(err)
	}
	err = c.metaFile.Close()
	if err != nil {
		panic(err)
	}
}
