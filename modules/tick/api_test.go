package tick_test

import (
	"testing"

	"github.com/kroppt/voxels/modules/camera"
	"github.com/kroppt/voxels/modules/tick"
)

func TestNewTickNotNil(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{})
	if tickMod == nil {
		t.Fatal("new tick module was nil")
	}
}

func TestNewTickPlayerNotNil(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic, but didn't")
		}
	}()
	tick.New(nil)
}

func TestGetCurrentTick(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{})
	expected := 0
	actual := tickMod.GetTick()

	if actual != expected {
		t.Fatalf("expected to get tick %v but got %v", expected, actual)
	}
}

func TestAdvanceTick(t *testing.T) {
	t.Parallel()
	tickMod := tick.New(&camera.FnModule{})
	expected := 1
	tickMod.AdvanceTick()
	actual := tickMod.GetTick()

	if actual != expected {
		t.Fatalf("expected to get tick %v but got %v", expected, actual)
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
	tickMod := tick.New(cameraMod)
	tickMod.AdvanceTick()
	if actual != expected {
		t.Fatal("expected camera to receive tick, but didn't")
	}
}
