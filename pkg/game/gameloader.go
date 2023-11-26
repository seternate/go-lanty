package game

import (
	"io/fs"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
)

func LoadFromDirectory(d string) (games Games, err error) {
	walkDirFunc := func(path string, dirEntry fs.DirEntry, errParam error) (errReturn error) {
		if errParam != nil {
			log.Debug().Err(errParam).Str("directory", path).Msg("failed to walk into directory")
			return nil
		}
		if !dirEntry.IsDir() && (filepath.Ext(dirEntry.Name()) == ".yaml" || filepath.Ext(dirEntry.Name()) == ".yml") {
			game := Game{}
			errReturn := filesystem.LoadFromYAMLFile(path, &game)
			if errReturn != nil {
				log.Warn().Err(errReturn).Str("file", path).Msg("error loading YAML file of game")
				return nil
			}
			errReturn = game.ValidateLazy()
			if errReturn != nil {
				log.Warn().Err(errReturn).Str("file", path).Msg("error on validation while loading game from YAML file")
				return nil
			}
			game.NormalizeState()
			log.Trace().Str("file", path).Msg("loaded game config YAML file")
			errReturn = games.Add(game)
			if errReturn != nil {
				log.Error().Err(errReturn).Str("slug", game.Slug).Msg("error adding game to games")
				return nil
			}
			log.Trace().Str("slug", game.Slug).Msg("added game to games")
		}
		return nil
	}
	err = filepath.WalkDir(d, walkDirFunc)
	return
}
