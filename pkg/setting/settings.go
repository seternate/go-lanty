package setting

import "github.com/seternate/go-lanty/pkg/filesystem"

const (
	SETTINGS_PATH = "settings.yaml"
)

type Settings struct {
	ServerPort          int    `yaml:"serverport"`
	GameConfigDirectory string `yaml:"game-config-directory"`
	GameFileDirectory   string `yaml:"game-file-directory"`
	GameIconDirectory   string `yaml:"game-icon-directory"`
}

func LoadSettings() (settings Settings, err error) {
	err = filesystem.LoadFromYAMLFile(SETTINGS_PATH, &settings)
	return
}
