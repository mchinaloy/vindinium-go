package main

import (
	"fmt"

	"github.com/mikechinaloy/vindinium-go/bot"
	"github.com/mikechinaloy/vindinium-go/request"
)

const (
	arenaURL    = "http://vindinium.org/api/arena"
	trainingURL = "http://vindinium.org/api/training"
)

func main() {
	training()
}

func arena() {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	play(arenaURL, body)
}

func training() {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	body["map"] = "m6"
	play(trainingURL, body)
}

func play(url string, body map[string]string) {
	gameState := request.PostRequest(url, body)
	fmt.Printf("\nGame starting in mode: %s", url)
	for gameState.Game.Finished != true && gameState.Hero.Crashed != true {
		gameState = bot.Move(gameState.PlayURL, gameState)
		fmt.Printf("\nMove #: %d of %d", gameState.Game.Turn, gameState.Game.MaxTurns)
	}
	fmt.Printf("\nGame finished, view the replay here: %s", gameState.ViewURL)
}
