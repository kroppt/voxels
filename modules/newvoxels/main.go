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
	"github.com/kroppt/voxels/modules/tick"
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
	if readCloser, err := fileMod.GetReadCloser("settings.conf"); err != nil {
		log.Warn(err)
	} else {
		settingsRepo.SetFromReader(readCloser)
		readCloser.Close()
	}

	graphicsMod := graphics.New(settingsRepo)
	err := graphicsMod.CreateWindow("newvoxels")
	if err != nil {
		log.Fatal(err)
	}

	testGen := &world.FnGenerator{
		FnGenerateChunk: func(key chunk.Position) chunk.Chunk {
			newChunk := chunk.New(key, settingsRepo.GetChunkSize())
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
	worldMod := world.New(graphicsMod, testGen, settingsRepo)
	playerMod := player.New(worldMod, settingsRepo, graphicsMod)
	cameraMod := camera.New(playerMod, player.PositionEvent{X: 0.5, Y: 0.5, Z: 3})
	inputMod := input.New(graphicsMod, cameraMod, settingsRepo)
	tickRateNano := int64(100 * 1e6)
	tickMod := tick.New(cameraMod, tick.FnTime{}, tickRateNano)
	graphicsMod.ShowWindow()

	keepRunning := true
	for keepRunning {
		if tickMod.IsNextTickReady() {
			tickMod.AdvanceTick()
		}
		graphicsMod.Render()
		keepRunning = inputMod.RouteEvents()
	}
	util.LogMetrics()
}
