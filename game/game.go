package game

import (
	"time"
)

type Game struct {
	now      func() time.Time
	tick     int
	lastTick time.Time
	thisTick time.Time
}

func New(now func() time.Time) *Game {
	currTime := now()
	return &Game{
		now:      now,
		tick:     0,
		lastTick: currTime,
		thisTick: currTime,
	}
}

func (g *Game) GetTick() int {
	return g.tick
}

func (g *Game) NextTick() {
	g.lastTick = g.thisTick
	g.thisTick = g.now()
	g.tick++
}

func (g *Game) TickDuration() time.Duration {
	return g.thisTick.Sub(g.lastTick)
}
