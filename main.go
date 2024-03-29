package main

import (
	"os"
	"sync"
	"time"

	"github.com/kroppt/voxels/log"
	"github.com/kroppt/voxels/modules/cache"
	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/file"
	"github.com/kroppt/voxels/modules/graphics"
	"github.com/kroppt/voxels/modules/input"
	"github.com/kroppt/voxels/modules/player"
	"github.com/kroppt/voxels/modules/tick"
	"github.com/kroppt/voxels/modules/view"
	"github.com/kroppt/voxels/modules/world"
	"github.com/kroppt/voxels/repositories/settings"
	"github.com/kroppt/voxels/util"
	"github.com/spf13/afero"
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

	graphicsMod := graphics.NewParallel(settingsRepo)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		graphicsMod.Run()
		wg.Done()
	}()

	err := graphicsMod.CreateWindow("voxels")
	if err != nil {
		log.Fatal(err)
	}
	generator := world.NewAlexWorldGenerator(settingsRepo)
	// generator := world.NewFlatWorldGenerator(settingsRepo)
	cacheMod := cache.New(afero.NewOsFs(), settingsRepo)
	viewMod := view.NewParallel(graphicsMod, settingsRepo)
	wg.Add(1)
	go func() {
		viewMod.Run()
		wg.Done()
	}()
	worldMod := world.NewParallel(graphicsMod, generator, settingsRepo, cacheMod, viewMod)
	wg.Add(1)
	go func() {
		worldMod.Run()
		wg.Done()
	}()
	playerMod := player.New(worldMod, settingsRepo, viewMod)
	cameraMod := camera.New(playerMod, player.PositionEvent{X: 0.5, Y: 20, Z: 0.5})
	inputMod := input.New(graphicsMod, cameraMod, settingsRepo, playerMod)
	tickRateNano := int64(1 * 1e6)
	tickMod := tick.New(cameraMod, tick.FnTime{}, tickRateNano)
	graphicsMod.ShowWindow()

	keepRunning := true
	before := time.Now()
	frames := 0
	for keepRunning {
		if tickMod.IsNextTickReady() {
			tickMod.AdvanceTick()
		}
		graphicsMod.Render()
		frames++
		keepRunning = inputMod.RouteEvents()
	}
	duration := time.Since(before)
	log.Perff("frames: %v, duration: %v, fps: %v", frames, duration, float64(frames)/duration.Seconds())
	worldMod.Close()
	wg.Wait()
	util.LogMetrics()
}
