package game

import (
	"os/exec"
	"path/filepath"

	"github.com/seternate/go-lanty-server/pkg/filesystem"
)

type Game struct {
	Slug             string `json:"slug" yaml:"slug"`
	Name             string `json:"name" yaml:"name"`
	ClientExecutable string `json:"clientexecutable" yaml:"clientexecutable"`
	ServerExecutable string `json:"serverexecutable" yaml:"serverexecutable"`
}

type Games []Game

func (game *Game) getAbsolutePathToExecutable(gameDirectory string) (string, error) {
	gameExecutable := filepath.Base(game.ClientExecutable)
	gamePath := filesystem.SearchFileByName(gameExecutable, gameDirectory)[0]

	return filepath.Abs(gamePath)
}

func (game *Game) Start(gameDirectory string) (*exec.Cmd, error) {
	gamePath, _ := game.getAbsolutePathToExecutable(gameDirectory)

	cmd := exec.Command("./" + filepath.Base(game.ClientExecutable))
	cmd.Dir = filepath.Dir(gamePath)
	err := cmd.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (game *Game) OpenFile(gameDirectory string) error {
	gamePath, _ := game.getAbsolutePathToExecutable(gameDirectory)

	cmd := exec.Command("explorer", "/select,", gamePath)
	return cmd.Run()
}
