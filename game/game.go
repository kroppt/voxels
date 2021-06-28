package game

import (
	"time"
)

// Game handles high-level details of the game.
type Game struct {
	now      func() time.Time
	tick     int
	lastTick time.Time
	thisTick time.Time
}

// New returns a new Game.
func New(now func() time.Time) *Game {
	currTime := now()
	return &Game{
		now:      now,
		tick:     0,
		lastTick: currTime,
		thisTick: currTime,
	}
}

// GetTick gets the current tick count.
func (g *Game) GetTick() int {
	return g.tick
}

// NextTick iterates the game to the next tick.
func (g *Game) NextTick() {
	g.lastTick = g.thisTick
	g.thisTick = g.now()
	g.tick++
}

// GetTickDuration returns the duration the current tick is responsible for.
func (g *Game) GetTickDuration() time.Duration {
	return g.thisTick.Sub(g.lastTick)
}

// OsTimeNow returns the current time according to the operating system.
func OsTimeNow() time.Time {
	return time.Now()
}
