package main

import (
	"os"

	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/events"
	"github.com/kroppt/voxels/modules/graphics"
)

func main() {
	log.SetInfoOutput(os.Stderr)
	log.SetWarnOutput(os.Stderr)
	log.SetDebugOutput(os.Stderr)
	log.SetPerfOutput(os.Stderr)
	log.SetFatalOutput(os.Stderr)
	log.SetColorized(false)

	graphicsMod := graphics.New()
	eventsMod := events.New(graphicsMod)

	err := graphicsMod.CreateWindow("newvoxels", 1920, 1080)
	if err != nil {
		log.Fatal(err)
	}

	graphicsMod.ShowWindow()

	eventsMod.RouteEvents()
}
