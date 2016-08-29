package aibot

import (
	"github.com/mikechinaloy/vindinium-go/model"
	"github.com/mikechinaloy/vindinium-go/request"
)

// SamuraiBot fights
type SamuraiBot struct {
}

// Move murders
func (i *SamuraiBot) Move(url string, gameState *model.GameState) *model.GameState {
	board := model.ParseBoard(&gameState.Game.Board)
	newGameState := request.PostRequest(url, makeDecision(board))
	return newGameState
}

func makeDecision(board *[][]model.Coordinate) map[string]string {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	return body
}