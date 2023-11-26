package game

import (
	"io/fs"
	"path/filepath"

	"github.com/seternate/go-lanty-server/pkg/filesystem"
)

func LoadFromDirectory(d string) (games Games, err error) {
	walkDirFunc := func(path string, dirEntry fs.DirEntry, errParam error) (errReturn error) {
		if errParam != nil {
			return nil
		}

		if !dirEntry.IsDir() && (filepath.Ext(dirEntry.Name()) == ".yaml" || filepath.Ext(dirEntry.Name()) == ".yml") {
			game := Game{}
			errReturn := filesystem.LoadFromYAMLFile(path, &game)
			if errReturn != nil {
				return errReturn
			}
			games = append(games, game)
		}

		return nil
	}

	filepath.WalkDir(d, walkDirFunc)

	return
}
