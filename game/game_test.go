package game_test

import (
	"testing"
	"time"

	"github.com/kroppt/voxels/game"
)

func osTimeNow() time.Time {
	return time.Now()
}

func TestNewGame(t *testing.T) {

	t.Run("should not return nil", func(t *testing.T) {
		t.Parallel()

		g := game.New(osTimeNow)
		if g == nil {
			t.Fatal("expect non-nil pointer but got nil")
		}
	})

}

func TestGetTick(t *testing.T) {

	t.Run("should return 0 for a new game", func(t *testing.T) {
		t.Parallel()

		g := game.New(osTimeNow)
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

		g := game.New(osTimeNow)
		g.NextTick()
		expect := 1

		actual := g.GetTick()

		if actual != expect {
			t.Fatalf("expected %v but got %v", expect, actual)
		}
	})

}

func TestTickDuration(t *testing.T) {

	t.Run("first tick should have 0 duration", func(t *testing.T) {
		t.Parallel()

		g := game.New(osTimeNow)
		expect := time.Duration(0)

		actual := g.TickDuration()

		if actual != expect {
			t.Fatalf("expected %v but got %v", expect, actual)
		}
	})

	t.Run("second tick should be time since first tick", func(t *testing.T) {
		now := time.Now()
		stubTimeNow := func() time.Time {
			return now
		}
		g := game.New(stubTimeNow)
		expect := time.Second

		now = now.Add(time.Second)
		g.NextTick()
		actual := g.TickDuration()

		if actual != expect {
			t.Fatalf("expected %v but got %v", expect, actual)
		}
	})

}
