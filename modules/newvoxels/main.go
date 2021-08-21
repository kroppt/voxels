package main

import (
	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/events"
	"github.com/kroppt/voxels/modules/graphics"
)

func main() {
	graphicsMod := graphics.New()
	eventsMod := events.New(graphicsMod)

	err := graphicsMod.CreateWindow("newvoxels", 1920, 1080)
	if err != nil {
		log.Fatal(err)
	}

	graphicsMod.ShowWindow()

	eventsMod.RouteEvents()
}
