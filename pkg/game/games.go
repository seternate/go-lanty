package game

import (
	"errors"
	"slices"
)

type Games struct {
	games map[string]Game
	slugs []string
}

func (games *Games) HasGame(slug string) bool {
	return slices.Contains(games.Slugs(), slug)
}

func (games *Games) Add(game Game) (err error) {
	if len(games.slugs) == 0 {
		games.games = make(map[string]Game)
	} else if games.HasGame(game.Slug) {
		return errors.New("game already in list")
	}
	games.slugs = append(games.slugs, game.Slug)
	games.games[game.Slug] = game
	return
}

func (games *Games) Equal(g Games) bool {
	equalSlugs := slices.Compare(games.Slugs(), g.Slugs())
	if equalSlugs != 0 {
		return false
	}
	for _, slug := range games.Slugs() {
		left, _ := games.Get(slug)
		right, _ := g.Get(slug)
		if left.Equal(right) {
			return true
		}
	}
	return false
}

func (games *Games) Get(slug string) (game Game, err error) {
	game, found := games.games[slug]
	if !found {
		err = errors.New("no game with specified slug found")
	}
	return
}

func (games *Games) Size() int {
	return len(games.slugs)
}

func (games *Games) Slugs() []string {
	return games.slugs
}

func (games Games) Games() (gamelist []Game) {
	for _, slug := range games.slugs {
		game, _ := games.Get(slug)
		gamelist = append(gamelist, game)
	}
	return
}
