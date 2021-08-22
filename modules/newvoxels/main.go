package main

import (
	"os"

	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/chunk"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/modules/player"
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
	playerMod := player.New(chunkMod)
	inputMod := input.New(graphicsMod, playerMod)

	err := graphicsMod.CreateWindow("newvoxels", 1920, 1080)
	if err != nil {
		log.Fatal(err)
	}

	graphicsMod.ShowWindow()

	inputMod.RouteEvents()
}
