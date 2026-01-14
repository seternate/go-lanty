package game

type GameRepository interface {
	GetGame(slug string) (*Game, error)
	SaveGame(*Game) error
	DeleteGame(slug string) error
}
