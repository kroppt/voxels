package main

import (
	"os"

	"github.com/kroppt/voxels/chunk"
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/file"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/kroppt/voxels/util"
)

func main() {
	log.SetInfoOutput(os.Stderr)
	log.SetWarnOutput(os.Stderr)
	log.SetDebugOutput(os.Stderr)
	log.SetPerfOutput(os.Stderr)
	log.SetFatalOutput(os.Stderr)
	log.SetColorized(false)
	util.SetMetricsEnabled(true)

	fileMod := file.New()
	settingsRepo := settings.New()
	graphicsMod := graphics.New(settingsRepo)
	chunkSize := uint32(1)
	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.Position) chunk.Chunk {
			newChunk := chunk.New(key, chunkSize)
			if key == (chunk.Position{X: 0, Y: 0, Z: 0}) {
				newChunk.SetBlockType(chunk.VoxelCoordinate{
					X: 0,
					Y: 0,
					Z: 0,
				}, chunk.BlockTypeDirt)
			}
			return newChunk
		},
	}
	worldMod := world.New(graphicsMod, testGen)
	playerMod := player.New(worldMod, settingsRepo, graphicsMod, chunkSize)
	cameraMod := camera.New(playerMod, player.PositionEvent{})
	inputMod := input.New(graphicsMod, cameraMod, settingsRepo)

	if readCloser, err := fileMod.GetReadCloser("settings.conf"); err != nil {
		log.Warn(err)
	} else {
		settingsRepo.SetFromReader(readCloser)
		readCloser.Close()
	}

	err := graphicsMod.CreateWindow("newvoxels")
	if err != nil {
		log.Fatal(err)
	}

	graphicsMod.ShowWindow()

	keepRunning := true
	for keepRunning {
		graphicsMod.Render()
		keepRunning = inputMod.RouteEvents()
	}
	util.LogMetrics()
}
