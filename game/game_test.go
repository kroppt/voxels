package game_test

import (
	"testing"

	"github.com/kroppt/voxels/game"
)

func TestNewGame(t *testing.T) {

	t.Run("should not return nil", func(t *testing.T) {
		t.Parallel()

		g := game.New()
		if g == nil {
			t.Fatal("expect non-nil pointer but got nil")
		}
	})

}

func TestGetTick(t *testing.T) {

	t.Run("should return 0 for a new game", func(t *testing.T) {
		t.Parallel()

		g := game.New()
		expect := 0

		actual := g.GetTick()

		if actual != expect {
			t.Fatalf("expected %v but got %v", expect, actual)
		}
	})

}

func TestNextTick(t *testing.T) {

	t.Run("should increment tick 0 to 1", func(t *testing.T) {
		t.Parallel()

		g := game.New()
		g.NextTick()
		expect := 1

		actual := g.GetTick()

		if actual != expect {
			t.Fatalf("expected %v but got %v", expect, actual)
		}
	})

}
