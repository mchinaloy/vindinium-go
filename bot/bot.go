package bot

import (
	"math/rand"

	"github.com/mikechinaloy/vindinium-go/model"
	"github.com/mikechinaloy/vindinium-go/request"
)

// Move makes a random decision to move
func Move(url string, gameState model.GameState) model.GameState {
	r := rand.Int31n(12)
	var newGameState model.GameState
	if r <= 2 {
		newGameState = request.PostRequest(url, createBody(model.Stay))
	} else if r <= 4 {
		newGameState = request.PostRequest(url, createBody(model.North))
	} else if r <= 6 {
		newGameState = request.PostRequest(url, createBody(model.South))
	} else if r <= 8 {
		newGameState = request.PostRequest(url, createBody(model.East))
	} else {
		newGameState = request.PostRequest(url, createBody(model.West))
	}
	return newGameState
}

func createBody(direction string) map[string]string {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	body["dir"] = direction
	return body
}
