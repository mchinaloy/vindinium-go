package ai

import "github.com/mikechinaloy/vindinium-go/model"

// Bot defines generic bot behaviours
type Bot interface {
	Move(url string, gameState *model.GameState, key string) *model.GameState
}
