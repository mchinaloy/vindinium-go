package model

// GameState represents current game
type GameState struct {
	id       string
	turn     int
	maxTurns int
	heroes   []Hero
	token    string
	viewURL  string
	playURL  string
	board    Board
}
