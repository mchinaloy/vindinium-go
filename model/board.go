package model

import "fmt"

// Board represents the tiles in the game
type Board struct {
	Size     int
	Tiles    string
	Finished bool
}

func ParseBoard(board *Board) *[][]Coordinate {
	x := board.Size
	y := board.Size

	tiles := make([][]Coordinate, x)

	for i := range tiles {
		tiles[i] = make([]Coordinate, y)
	}

	startIndex := 0
	multiplier := 2
	fmt.Println("Board size: ", board.Size)
	for i := 0; i < board.Size; i++ {
		boardLine := board.Tiles[startIndex:(board.Size) * multiplier]
		fmt.Println(boardLine)
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

	return &tiles
}

func createCoordinate (boardTile string, row int, column int) *Coordinate {
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
		coordinate.Type = Mine
		return coordinate
	case boardTile[0:1] == "$" :
		coordinate.Type = PlayerMine
		return coordinate
	default:
		coordinate.Type = Player
		return coordinate
	}
}
