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
	file         afero.File
	settingsRepo settings.Interface
}

func (c *core) save(ch chunk.Chunk) {
	curr := int64(0)
	for {
		chunkSize := c.settingsRepo.GetChunkSize()
		numElems := chunk.VertSize * chunkSize * chunkSize * chunkSize
		byteSize := chunk.BytesPerElement * numElems
		bs := make([]byte, byteSize)
		n, err := c.file.ReadAt(bs, curr)
		if n == 0 && errors.Is(err, io.EOF) {
			break
		}
		if uint32(n) != byteSize {
			log.Printf("expected %v bytes to be read, but only %v were read", byteSize, n)
			return
		}
		if !errors.Is(err, io.EOF) && err != nil {
			log.Print(err)
			return
		}
		buf := bytes.NewBuffer(bs)
		flatData := make([]float32, numElems)
		err = binary.Read(buf, binary.LittleEndian, flatData)
		if err != nil {
			log.Print(err)
			return
		}
		firstVox := chunk.VoxelCoordinate{
			X: int32(flatData[0]),
			Y: int32(flatData[1]),
			Z: int32(flatData[2]),
		}
		key := chunk.VoxelCoordToChunkCoord(firstVox, chunkSize)
		if key == ch.Position() {
			// overwrite existing chunk in file with update info
			var buf bytes.Buffer
			err := binary.Write(&buf, binary.LittleEndian, ch.GetFlatData())
			if err != nil {
				log.Print(err)
				return
			}
			n, err := c.file.WriteAt(buf.Bytes(), curr)
			expectSize := buf.Len()
			if n != expectSize {
				log.Printf("expected to write %v bytes, but only wrote %v bytes", expectSize, n)
				return
			}
			if err != nil {
				log.Print(err)
				return
			}
			return
		}
		curr += int64(byteSize)
	}
	// one wasn't overwritten, append new chunk to end

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, ch.GetFlatData())
	if err != nil {
		log.Print(err)
		return
	}
	n, err := c.file.WriteAt(buf.Bytes(), curr)
	expectSize := buf.Len()
	if n != expectSize {
		log.Printf("expected to write %v bytes, but only wrote %v bytes", expectSize, n)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
	return
}

func (c *core) load(pos chunk.Position) (chunk.Chunk, bool) {
	curr := int64(0)
	for {
		chunkSize := c.settingsRepo.GetChunkSize()
		numElems := chunk.VertSize * chunkSize * chunkSize * chunkSize
		byteSize := chunk.BytesPerElement * numElems
		bs := make([]byte, byteSize)
		n, err := c.file.ReadAt(bs, curr)
		if n == 0 && errors.Is(err, io.EOF) {
			log.Printf("did not find chunk %v in file", pos)
			return chunk.Chunk{}, false
		}
		if uint32(n) != byteSize {
			log.Printf("expected %v bytes to be read, but only %v were read", byteSize, n)
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
		firstVox := chunk.VoxelCoordinate{
			X: int32(flatData[0]),
			Y: int32(flatData[1]),
			Z: int32(flatData[2]),
		}
		key := chunk.VoxelCoordToChunkCoord(firstVox, chunkSize)
		if key == pos {
			return chunk.NewFromData(flatData, chunkSize, key), true
		}
		curr += int64(byteSize)
	}
}
