package aibot

import (
	"github.com/mikechinaloy/vindinium-go/model"
	"github.com/mikechinaloy/vindinium-go/request"
	"github.com/nickdavies/go-astar/astar"
	"fmt"
	"math/rand"
	"strings"
	"strconv"
)

const (
	Direction = "dir"
)

// AStarBot
type AStarBot struct {
}

var lastMoveDir string
var lastMoveCount int
var goingAround bool

// Move murders
func (i *AStarBot) Move(url string, gameState *model.GameState, key string) *model.GameState {
	board := model.ParseBoard(&gameState.Game.Board)
	newGameState := request.PostRequest(url, makeDecision(board, gameState, key))
	return newGameState
}

func makeDecision(board [][]model.Coordinate, gameState *model.GameState, key string) map[string]string {
	var body = make(map[string]string)
	body["key"] = key

	var target []astar.Point

	if gameState.Hero.Life < 10 {
		fmt.Printf("\nHealth is below 10 at: %d, looking for a drink", gameState.Hero.Life)
		target = findTavern(gameState, board)
	} else {
		target = findGoldMine(gameState, board)
	}

	route := findShortestPath(board, gameState, target)

	if route != nil {
		body[Direction] = getDirectionFromStartTarget([]astar.Point{{Row: gameState.Hero.Pos.X, Col: gameState.Hero.Pos.Y}}, []astar.Point{{Row: route.Parent.Row, Col: route.Parent.Col}})
	} else {
		fmt.Println("\nNo route could be determined, randomizing move")
		r := rand.Int31n(12)
		if r <= 2 {
			body[Direction] = model.Stay
		} else if r <= 4 {
			body[Direction] = model.North
		} else if r <= 6 {
			body[Direction] = model.South
		} else if r <= 8 {
			body[Direction] = model.East
		} else {
			body[Direction] = model.West
		}
	}

	if lastMoveDir == body[Direction] {
		lastMoveCount++
	} else {
		lastMoveDir = body[Direction]
		lastMoveCount = 0
	}

	if lastMoveCount >= gameState.Game.Board.Size / 2 {
		fmt.Printf("\nStuck at x: %d, y%d, attempting to go around", gameState.Hero.Pos.X, gameState.Hero.Pos.Y)
		goingAround = true
		if gameState.Hero.Pos.X - 1 < 0 {
			if gameState.Hero.Pos.Y - 1 < 0 {
				body[Direction] = model.East
			} else {
				body[Direction] = model.West
			}
		} else if gameState.Hero.Pos.Y + 1 > gameState.Game.Board.Size {
			if gameState.Hero.Pos.X + 1 < 0 {
				body[Direction] = model.East
			} else {
				body[Direction] = model.West
			}
		} else {
			if gameState.Hero.Pos.Y - 1 < 0 {
				body[Direction] = model.East
			} else {
				body[Direction] = model.West
			}
		}
	}

	if goingAround && (lastMoveDir == model.East || lastMoveDir == model.West) {
		body[Direction] = model.South
		goingAround = false
	}

	fmt.Printf("\nMoving in direction: %s", body[Direction])

	return body
}

func getDirectionFromStartTarget(start []astar.Point, target []astar.Point) string {
	if start[0].Row > target[0].Row {
		return model.North
	} else if start[0].Row < target[0].Row {
		return model.South
	} else if start[0].Col > target[0].Col {
		return model.West
	} else {
		return model.East
	}
}

func findTavern(gameState *model.GameState, board [][]model.Coordinate) []astar.Point {
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if board[i][j].Type == model.Tavern {
				fmt.Printf("\nLocated a tavern to go after x: %d, y: %d", i, j)
				return []astar.Point{{Row: i, Col: j}}
			}
		}
	}
	return nil;
}

func findGoldMine(gameState *model.GameState, board [][]model.Coordinate) []astar.Point {
	var playerId = model.PlayerMine + strconv.Itoa(gameState.Hero.Id);
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if strings.Contains(board[i][j].Type, model.PlayerMine) {
				if board[i][j].Type == playerId {
					fmt.Printf("\nThe closest mine is owned by me x: %d, y: %d, looking for an alternative", i, j)
					continue
				}
			} else {
				continue
			}
			fmt.Printf("\nLocated a player mine to go after x: %d, y: %d", i, j)
			return []astar.Point{{Row: i, Col: j}}
		}
	}
	return nil;
}

func findShortestPath(board [][] model.Coordinate, gameState *model.GameState, target []astar.Point) *astar.PathPoint {
	// Create blank A* graph
	graph := astar.NewAStar(gameState.Game.Board.Size, gameState.Game.Board.Size)
	p2p := astar.NewPointToPoint()

	// Fill in obstacles
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if board[j][i].Type == model.Wood {
				graph.FillTile(astar.Point{j, i}, -1)
			} else if strings.Contains(board[j][i].Type, model.PlayerMine) {
				graph.FillTile(astar.Point{j, i}, 5)
			} else if board[j][i].Type == model.Tavern {
				graph.FillTile(astar.Point{j, i}, 5)
			}
		}
	}
	source := []astar.Point{{Row: gameState.Hero.Pos.X, Col: gameState.Hero.Pos.Y}}
	return graph.FindPath(p2p, source, target)
}