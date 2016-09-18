package aibot

import (
	"github.com/mikechinaloy/vindinium-go/model"
	"github.com/mikechinaloy/vindinium-go/request"
	"github.com/nickdavies/go-astar/astar"
	"math/rand"
	"strings"
	"strconv"
	"fmt"
)

const (
	Direction = "dir"
)

// AStarBot
type AStarBot struct {
}

var (
	lastMoveDir string
	lastMoveCount int
	randomCount int
	murderMode bool
	campingAtTavern bool
)

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

	player := shouldSwitchToMurderMode(gameState)

	gameTurnProgress := (gameState.Game.Turn / gameState.Game.MaxTurns) * 100

	nearestTavern := findNearestObjective(gameState, board, model.Tavern, false)

	if isAtTavern(gameState, nearestTavern) && gameState.Hero.Gold > 0 && gameState.Hero.Life < 100 {
		target = nearestTavern
		fmt.Printf("\n At a tavern at=%s, drinking to get full HP", target)
	} else {
		if gameState.Hero.Life < 60 || shouldStandStillAtTavern(gameState) && gameState.Hero.Gold > 0 {
			fmt.Printf("\n HP is=%d, searching for a tavern", gameState.Hero.Life)
			target = findNearestObjective(gameState, board, model.Tavern, false)
		} else if murderMode {
			var playerId = model.Player + strconv.Itoa(player.Id);
			target = findObjective(gameState, board, playerId)
			fmt.Printf("\n Player=%s holds a lot of mines, going after them at=%s", playerId, target)
		} else {
			if existsThereNeutralMines(gameState, board) && gameTurnProgress > 60 {
				target = findNearestObjective(gameState, board, model.NeutralMine, false)
				fmt.Println("Going for a Neutral Mine at=", target)
			} else {
				target = findNearestObjective(gameState, board, model.PlayerMine, true)
				fmt.Println("Going for a Player Mine at=", target)
			}
		}
	}

	route := findShortestPath(board, gameState, target)

	if route != nil {
		body[Direction] = getDirectionFromStartTarget([]astar.Point{{Row: gameState.Hero.Pos.X, Col: gameState.Hero.Pos.Y}}, []astar.Point{{Row: route.Parent.Row, Col: route.Parent.Col}})
	} else {
		fmt.Println("\nFailure to determine shortest route, randomizing move")
		body[Direction] = makeRandomMove()
	}

	if lastMoveDir == body[Direction] {
		lastMoveCount++
	} else {
		lastMoveDir = body[Direction]
		lastMoveCount = 0
	}

	if (lastMoveCount >= (gameState.Game.Board.Size / 2) + 1 && randomCount < 2) && !campingAtTavern {
		fmt.Printf("\nProbably Stuck at x=%d, y=%d, attempting to random twice", gameState.Hero.Pos.X, gameState.Hero.Pos.Y)
		body[Direction] = makeRandomMove()
	}

	if (randomCount == 1) {
		randomCount = 0;
	}

	fmt.Printf("\nMoving in direction=%s", body[Direction])

	return body
}

func isAtTavern(gameState *model.GameState, target []astar.Point) bool {
	if gameState.Hero.Pos.X == target[0].Row && gameState.Hero.Pos.Y == target[0].Col {
		return true
	}
	return false
}

func shouldStandStillAtTavern(gameState *model.GameState) bool {
	playerWithMostMines := gameState.Game.Heroes[0]

	for i := 0; i < len(gameState.Game.Heroes); i ++ {
		if gameState.Game.Heroes[i].MineCount > playerWithMostMines.MineCount {
			playerWithMostMines = gameState.Game.Heroes[i]
		}
	}

	if playerWithMostMines.MineCount > 0 {
		percentageOfMines := (float64(playerWithMostMines.MineCount) / float64(model.TotalNumberOfMines)) * 100
		if (percentageOfMines >= 50 && playerWithMostMines.Id == gameState.Hero.Id) || gameState.Hero.MineCount == model.TotalNumberOfMines {
			campingAtTavern = true
			fmt.Println("Camping at the nearest tavern because we are doing well!")
			return true
		}
	}

	campingAtTavern = false
	return false
}

func makeRandomMove() string {
	randomCount++
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

func existsThereNeutralMines(gameState *model.GameState, board [][]model.Coordinate) bool {
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if board[i][j].Type == model.NeutralMine {
				fmt.Printf("\nFound a neutral mine at x=%d, y=%d", i, j)
				return true
			}
		}
	}
	return false
}

func shouldSwitchToMurderMode(gameState *model.GameState) *model.Hero {
	playerWithMostMines := gameState.Game.Heroes[0]

	for i := 0; i < len(gameState.Game.Heroes); i ++ {
		if gameState.Game.Heroes[i].MineCount > playerWithMostMines.MineCount {
			playerWithMostMines = gameState.Game.Heroes[i]
		}
	}

	if playerWithMostMines.MineCount > 0 {
		percentageOfMines := (float64(playerWithMostMines.MineCount) / float64(model.TotalNumberOfMines)) * 100
		fmt.Printf("TotalMines=%d, playerWithMostMines=%d, percentage of mines=%f", model.TotalNumberOfMines, playerWithMostMines.MineCount, percentageOfMines)
		if percentageOfMines >= 50 && playerWithMostMines.Life < gameState.Hero.Life && playerWithMostMines.Id != gameState.Hero.Id {
			murderMode = true
			return &playerWithMostMines
		}
	}
	murderMode = false
	return &playerWithMostMines
}

func findObjective(gameState *model.GameState, board [][]model.Coordinate, objective string) []astar.Point {
	row := 0
	col := 0
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if strings.Contains(board[i][j].Type, objective) {
				row = i
				col = j
			}
		}
	}
	return []astar.Point{{Row: row, Col: col}}
}

func findNearestObjective(gameState *model.GameState, board [][]model.Coordinate, objective string, isPlayerMine bool) []astar.Point {
	closestRow := 0
	closestCol := 0
	bestDiff := 0
	foundFirst := false

	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if strings.Contains(board[i][j].Type, objective) {
				fmt.Printf("\nLocated a %s to go after x=%d, y=%d", objective, i, j)
				if isPlayerMine {
					var playerId = model.PlayerMine + strconv.Itoa(gameState.Hero.Id);
					if strings.Contains(board[i][j].Type, model.PlayerMine) {
						if board[i][j].Type == playerId {
							fmt.Printf("\nFound a mine which is owned by me x=%d, y=%d, looking for an alternative", i, j)
							continue
						}
					}
				}
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

func calculateNearestDiff(gameState *model.GameState, row int, col int) (int, int) {
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

func findShortestPath(board [][] model.Coordinate, gameState *model.GameState, target []astar.Point) *astar.PathPoint {
	// Create blank A* graph
	graph := astar.NewAStar(gameState.Game.Board.Size, gameState.Game.Board.Size)
	p2p := astar.NewPointToPoint()

	// Fill in obstacles
	for i := 0; i < gameState.Game.Board.Size; i++ {
		for j := 0; j < gameState.Game.Board.Size; j++ {
			if board[j][i].Type == model.Wood {
				graph.FillTile(astar.Point{j, i}, -1)
			} else if board[j][i].Type == model.NeutralMine {
				graph.FillTile(astar.Point{j, i}, 1000)
			} else if strings.Contains(board[j][i].Type, model.PlayerMine) {
				graph.FillTile(astar.Point{j, i}, 1000)
			}
		}
	}
	source := []astar.Point{{Row: gameState.Hero.Pos.X, Col: gameState.Hero.Pos.Y}}
	return graph.FindPath(p2p, source, target)
}