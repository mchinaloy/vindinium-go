package model

import (
	"strings"
)

// Board represents the tiles in the game
type Board struct {
	Size     int
	Tiles    string
	Finished bool
}

var (
	TotalNumberOfMines int
)

func ParseBoard(board *Board) [][]Coordinate {

	TotalNumberOfMines = 0

	x := board.Size
	y := board.Size

	tiles := make([][]Coordinate, x)

	for i := range tiles {
		tiles[i] = make([]Coordinate, y)
	}

	startIndex := 0
	multiplier := 2
	for i := 0; i < board.Size; i++ {
		boardLine := board.Tiles[startIndex:(board.Size) * multiplier]
		lineStartIndex := 0
		lineOffset := 2
		for j := 0; j < board.Size; j++ {
			boardTile := boardLine[lineStartIndex: lineOffset]
			tiles[i][j] = *createCoordinate(boardTile, i, j)
			lineOffset = lineOffset + 2
			lineStartIndex = lineStartIndex + 2
		}

		startIndex = board.Size * multiplier
		multiplier = multiplier + 2
	}

	return tiles
}

func createCoordinate(boardTile string, row int, column int) *Coordinate {
	coordinate := new(Coordinate)
	coordinate.X = row
	coordinate.Y = column
	switch {
	case boardTile == "##" :
		coordinate.Type = Wood
		return coordinate
	case boardTile == "  " :
		coordinate.Type = Land
		return coordinate
	case boardTile == "[]" :
		coordinate.Type = Tavern
		return coordinate
	case boardTile == "$-" :
		coordinate.Type = NeutralMine
		TotalNumberOfMines++
		return coordinate
	case strings.Contains(boardTile, "$") :
		coordinate.Type = boardTile[0:2]
		TotalNumberOfMines++
		return coordinate
	default:
		coordinate.Type = boardTile[0:2]
		return coordinate
	}
}
