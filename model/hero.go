package model

// Hero represents a player
type Hero struct {
	id        int
	name      string
	userID    string
	elo       int
	pos       Position
	life      int
	gold      int
	mineCount int
	spawnPos  Position
	crashed   bool
}
