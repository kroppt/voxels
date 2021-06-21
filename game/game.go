package game

type Game struct {
	tick int
}

func New() *Game {
	return &Game{}
}

func (g *Game) GetTick() int {
	return g.tick
}

func (g *Game) NextTick() {
	g.tick++
}
