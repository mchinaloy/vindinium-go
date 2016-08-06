package model

// Hero represents a player
type Hero struct {
	Id        int
	Name      string
	UserId    string
	Elo       int
	Pos       Position
	LastDir   string
	Life      int
	Gold      int
	MineCount int
	SpawnPos  Position
	Crashed   bool
}
