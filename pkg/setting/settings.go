package setting

import "github.com/seternate/go-lanty/pkg/filesystem"

const (
	VERSION                   = "0.1.0-beta"
	SETTINGS_PATH             = "settings.yaml"
	CLIENT_DOWNLOAD_DIRECTORY = "download"
	CLIENT_DOWNLOAD_FILE      = "lanty.zip"
)

type Settings struct {
	ServerPort             int    `yaml:"serverport"`
	ServerGracefulShutdown int    `yaml:"-"`
	GameConfigDirectory    string `yaml:"game-config-directory"`
	GameFileDirectory      string `yaml:"game-file-directory"`
	GameIconDirectory      string `yaml:"game-icon-directory"`
}

func LoadSettings() (settings Settings, err error) {
	err = filesystem.LoadFromYAMLFile(SETTINGS_PATH, &settings)
	return
}
