package model

// Game represents the game
type Game struct {
	Id       string
	Turn     int
	MaxTurns int
	Heroes   []Hero
	Board    Board
	Finished bool
}
