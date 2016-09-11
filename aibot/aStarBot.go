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
var lastDrinkCount int
var goingAround bool
var attemptToGoAroundCount int

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

	if gameState.Hero.Life < 60 && lastDrinkCount < 5 {
		fmt.Printf("\nHealth is below 50 at: %d, looking for a drink for the %dth time", gameState.Hero.Life, lastDrinkCount)
		target = findNearestTavern(gameState, board)
		lastDrinkCount++
	} else {
		lastDrinkCount = 0
		target = findNearestGoldMine(gameState, board)
	}

	route := findShortestPath(board, gameState, target)

	if route != nil {
		body[Direction] = getDirectionFromStartTarget([]astar.Point{{Row: gameState.Hero.Pos.X, Col: gameState.Hero.Pos.Y}}, []astar.Point{{Row: route.Parent.Row, Col: route.Parent.Col}})
	} else {
		fmt.Println("\nNo route could be determined, randomizing move")
		body[Direction] = makeRandomMove()
	}

	if lastMoveDir == body[Direction] {
		lastMoveCount++
	} else {
		lastMoveDir = body[Direction]
		lastMoveCount = 0
	}

	if lastMoveCount >= (gameState.Game.Board.Size / 2) + 1 {
		fmt.Printf("\nProbably Stuck at x: %d, y%d, attempting to go around", gameState.Hero.Pos.X, gameState.Hero.Pos.Y)
		goingAround = true
		if(attemptToGoAroundCount > 5) {
			goingAround = false
			fmt.Printf("\nFailed to go around for %dth time, making a random move", attemptToGoAroundCount)
			body[Direction] = makeRandomMove()
			attemptToGoAroundCount = 0
		} else {
			body[Direction] = goAroundObstacle(gameState)
			attemptToGoAroundCount++
		}
	}

	if goingAround && (lastMoveDir == model.East || lastMoveDir == model.West) {
		body[Direction] = model.South
		goingAround = false
	}

	fmt.Printf("\nMoving in direction: %s", body[Direction])

	return body
}

func goAroundObstacle(gameState *model.GameState) string {
	if gameState.Hero.Pos.X - 1 < 0 {
		if gameState.Hero.Pos.Y - 1 < 0 {
			return model.East
		} else {
			return model.West
		}
	} else if gameState.Hero.Pos.Y + 1 > gameState.Game.Board.Size {
		if gameState.Hero.Pos.X + 1 < 0 {
			return model.East
		} else {
			return model.West
		}
	} else {
		if gameState.Hero.Pos.Y - 1 < 0 {
			return model.East
		} else {
			return model.West
		}
	}
}

func makeRandomMove() string {
	r := rand.Int31n(12)
	if r <= 2 {
		return model.Stay
	} else if r <= 4 {
		return model.North
	} else if r <= 6 {
		return model.South
	} else if r <= 8 {
		return model.East
	} else {
		return model.West
	}
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

func findNearestTavern(gameState *model.GameState, board [][]model.Coordinate) []astar.Point {
	closestRow := 0
	closestCol := 0
	bestDiff := 0
	foundFirst := false
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if board[i][j].Type == model.Tavern {
				fmt.Printf("\nLocated a tavern to go after x: %d, y: %d", i, j)
				if foundFirst == true {
					rowDiff, colDiff := calculateNearestDiff(gameState, i, j)
					if (rowDiff + colDiff) < bestDiff {
						bestDiff = rowDiff + colDiff
						closestRow = i
						closestCol = j
					}
				} else {
					rowDiff, colDiff := calculateNearestDiff(gameState, i, j)
					closestRow = i
					closestCol = j
					bestDiff = rowDiff + colDiff
					foundFirst = true
				}
			}
		}
	}
	return []astar.Point{{Row: closestRow, Col: closestCol}}
}

func calculateNearestDiff(gameState *model.GameState, row int, col int) (int, int){
	rowDiff := 0
	colDiff := 0
	if row > gameState.Hero.Pos.X {
		rowDiff = row - gameState.Hero.Pos.X
	} else {
		rowDiff = gameState.Hero.Pos.X - row
	}
	if col > gameState.Hero.Pos.Y {
		colDiff = col - gameState.Hero.Pos.Y
	} else {
		colDiff = gameState.Hero.Pos.Y - col
	}
	return rowDiff, colDiff
}

func findNearestGoldMine(gameState *model.GameState, board [][]model.Coordinate) []astar.Point {
	closestRow := 0
	closestCol := 0
	bestDiff := 0
	foundFirst := false
	var playerId = model.PlayerMine + strconv.Itoa(gameState.Hero.Id);
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if strings.Contains(board[i][j].Type, model.PlayerMine) {
				if board[i][j].Type == playerId {
					fmt.Printf("\nFound a mine which is owned by me x: %d, y: %d, looking for an alternative", i, j)
					continue
				}
			} else {
				continue
			}
			fmt.Printf("\nLocated a player mine to go after x: %d, y: %d, seeing if its the closest one", i, j)
			if foundFirst == true {
				rowDiff, colDiff := calculateNearestDiff(gameState, i, j)
				if (rowDiff + colDiff) < bestDiff {
					bestDiff = rowDiff + colDiff
					closestRow = i
					closestCol = j
				}
			} else {
				rowDiff, colDiff := calculateNearestDiff(gameState, i, j)
				closestRow = i
				closestCol = j
				bestDiff = rowDiff + colDiff
				foundFirst = true
			}
		}
	}
	return []astar.Point{{Row: closestRow, Col: closestCol}}
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