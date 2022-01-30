package tick_test

import (
	"testing"
	"time"

	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/tick"
)

func TestNewTickNotNil(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{}, tick.FnTime{}, 1)
	if tickMod == nil {
		t.Fatal("new tick module was nil")
	}
}

func TestNewTickCameraNotNil(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	tick.New(nil, tick.FnTime{}, 1)
}

func TestNewTickTimeModuleNotNil(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	tick.New(&camera.FnModule{}, nil, 1)
}

func TestTickRatePositive(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	tick.New(&camera.FnModule{}, tick.FnTime{}, 0)
}

func TestGetCurrentTick(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{}, tick.FnTime{}, 1)
	expected := 0
	actual := tickMod.GetTick()

	if actual != expected {
		t.Fatalf("expected to get tick %v but got %v", expected, actual)
	}
}

func TestAdvanceTick(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{}, tick.FnTime{}, 1)
	expected := 1
	tickMod.AdvanceTick()
	actual := tickMod.GetTick()

	if actual != expected {
		t.Fatalf("expected to get tick %v but got %v", expected, actual)
	}
}

func TestNextTickIsNotReadyIfNotEnoughTimeHasPassed(t *testing.T) {
	t.Parallel()
	expected := false
	current := time.Now()
	timeMod := tick.FnTime{
		FnNow: func() time.Time {
			return current
		},
	}
	tickMod := tick.New(&camera.Module{}, timeMod, 1)
	actual := tickMod.IsNextTickReady()

	if actual != expected {
		t.Fatal("expected next tick to not be ready, but it was")
	}
}

func TestNextTickIsReadyIfEnoughTimeHasPassed(t *testing.T) {
	t.Parallel()
	expected := true
	current := time.Now()
	timeMod := &tick.FnTime{
		FnNow: func() time.Time {
			return current
		},
	}
	tickRateNano := 500 * 1e6
	tickMod := tick.New(&camera.Module{}, timeMod, 500*1e6)
	timeMod.FnNow = func() time.Time {
		return current.Add(time.Duration(tickRateNano))
	}
	actual := tickMod.IsNextTickReady()

	if actual != expected {
		t.Fatal("expected next tick to be ready, but it was not")
	}
}

func TestNextTickIsNotReadyIfBarelyNotEnoughTimeHasPassed(t *testing.T) {
	t.Parallel()
	expected := false
	current := time.Now()
	timeMod := &tick.FnTime{
		FnNow: func() time.Time {
			return current
		},
	}
	tickRateNano := 500 * 1e6
	tickMod := tick.New(&camera.Module{}, timeMod, 500*1e6)
	timeMod.FnNow = func() time.Time {
		return current.Add(time.Duration(tickRateNano - 1))
	}
	actual := tickMod.IsNextTickReady()

	if actual != expected {
		t.Fatal("expected next tick to not be ready, but it was")
	}
}

func TestCameraReceivesTick(t *testing.T) {
	t.Parallel()
	expected := true
	actual := false
	cameraMod := &camera.FnModule{
		FnTick: func() {
			actual = true
		},
	}
	tickMod := tick.New(cameraMod, tick.FnTime{}, 1)
	tickMod.AdvanceTick()
	if actual != expected {
		t.Fatal("expected camera to receive tick, but didn't")
	}
}
