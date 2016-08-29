package aibot

import (
	"math/rand"

	"github.com/mikechinaloy/vindinium-go/model"
	"github.com/mikechinaloy/vindinium-go/request"
)

// RandomBot randoms
type RandomBot struct {
}

// Move makes a random decision to move
func (i *RandomBot) Move(url string, gameState *model.GameState) *model.GameState {
	r := rand.Int31n(12)
	if r <= 2 {
		return request.PostRequest(url, createBody(model.Stay))
	} else if r <= 4 {
		return request.PostRequest(url, createBody(model.North))
	} else if r <= 6 {
		return request.PostRequest(url, createBody(model.South))
	} else if r <= 8 {
		return request.PostRequest(url, createBody(model.East))
	} else {
		return request.PostRequest(url, createBody(model.West))
	}
}

func createBody(direction string) map[string]string {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	body["dir"] = direction
	return body
}
