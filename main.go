package main

import (
	"fmt"

	"github.com/mikechinaloy/vindinium-go/ai"
	"github.com/mikechinaloy/vindinium-go/aibot"
	"github.com/mikechinaloy/vindinium-go/request"
	"os"
)

const (
	arenaURL = "http://vindinium.org/api/arena"
	trainingURL = "http://vindinium.org/api/training"
)

var (
	key string
	url string
	mode string
)

func main() {
	mode = os.Args[1]
	key = os.Args[2]

	if len(os.Args) > 3 {
		url = os.Args[3]
	}

	fmt.Printf("\nStarting GoBot with settings mode=%s, key=%s", mode, key)

	selectMode(mode, new(aibot.AStarBot))
}

func selectMode(mode string, bot ai.Bot) {
	switch mode {
	case "arena" :
		arena(bot)
	default:
		training(bot)
	}
}

func arena(bot ai.Bot) {
	var body = make(map[string]string)
	body["key"] = key
	if url != "" {
		play(url, body, bot)
	} else {
		for {
			play(arenaURL, body, bot)
		}
	}
}

func training(bot ai.Bot) {
	var body = make(map[string]string)
	body["key"] = key
	body["map"] = "m3"
	body["turns"] = "200"
	if url != "" {
		play(url, body, bot)
	} else {
		play(trainingURL, body, bot)
	}
}

func play(url string, body map[string]string, bot ai.Bot) {
	gameState := request.PostRequest(url, body)
	fmt.Printf("body=%s", gameState.Game.Board)
	fmt.Printf("\n\nGame starting at url=%s, viewurl=%s\n", url, gameState.ViewURL)
	count := 0
	for gameState.Game.Finished != true && gameState.Hero.Crashed != true {
		gameState = bot.Move(gameState.PlayURL, gameState, key)
		count++
	}
	fmt.Printf("\nGame finished, view the replay here=%s", gameState.ViewURL)
}
