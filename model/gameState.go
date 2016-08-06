package model

// GameState represents current game
type GameState struct {
	Game    Game
	Hero    Hero
	Token   string
	ViewURL string
	PlayURL string
}
