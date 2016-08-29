package main

import (
	"fmt"

	"github.com/mikechinaloy/vindinium-go/ai"
	"github.com/mikechinaloy/vindinium-go/aibot"
	"github.com/mikechinaloy/vindinium-go/request"
	"os"
)

const (
	arenaURL    = "http://vindinium.org/api/arena"
	trainingURL = "http://vindinium.org/api/training"
)

func main() {
	mode := os.Args[1]
	botName := os.Args[2]

	fmt.Printf("\nStarting GoBot with settings: %s and %s\n", mode, botName)

	selectMode(mode, selectBot(botName))
}

func selectBot(botName string) ai.Bot {
	switch botName {
	case "samuraiBot" :
		return new(aibot.SamuraiBot)
	default:
		return new(aibot.RandomBot)
	}
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
	body["key"] = "awffitw4"
	play(arenaURL, body, bot)
}

func training(bot ai.Bot) {
	var body = make(map[string]string)
	body["key"] = "awffitw4"
	body["map"] = "m2"
	body["turns"] = "200"
	play(trainingURL, body, bot)
}

func play(url string, body map[string]string, bot ai.Bot) {
	gameState := request.PostRequest(url, body)
	fmt.Printf("\nGame starting in mode: %s", url)
	count := 0
	for gameState.Game.Finished != true && gameState.Hero.Crashed != true {
		gameState = bot.Move(gameState.PlayURL, gameState)
		fmt.Printf("\nMove #: %d of %d", gameState.Game.Turn, gameState.Game.MaxTurns)
		count++
	}
	fmt.Printf("\nGame finished, view the replay here: %s", gameState.ViewURL)
}
