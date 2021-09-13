package main

import (
	"os"

	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/file"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/repositories/settings"
)

func main() {
	log.SetInfoOutput(os.Stderr)
	log.SetWarnOutput(os.Stderr)
	log.SetDebugOutput(os.Stderr)
	log.SetPerfOutput(os.Stderr)
	log.SetFatalOutput(os.Stderr)
	log.SetColorized(false)

	graphicsMod := graphics.New()
	chunkMod := chunk.New()
	playerMod := player.New(chunkMod, graphicsMod)
	fileMod := file.New()
	settingsRepo := settings.New()
	inputMod := input.New(graphicsMod, playerMod, settingsRepo)

	if readCloser, err := fileMod.GetReadCloser("settings.conf"); err != nil {
		log.Warn(err)
	} else {
		settingsRepo.SetFromReader(readCloser)
		readCloser.Close()
	}
	width, height := settingsRepo.GetResolution()
	if width == 0 || height == 0 {
		width = 1280
		height = 720
	}
	err := graphicsMod.CreateWindow("newvoxels", width, height)
	if err != nil {
		log.Fatal(err)
	}

	graphicsMod.ShowWindow()

	inputMod.RouteEvents()
}
