package game

type Game struct {
	Slug             string `json:"slug" yaml:"slug"`
	Name             string `json:"name" yaml:"name"`
	ClientExecutable string `json:"clientexecutable" yaml:"clientexecutable"`
	ServerExecutable string `json:"serverexecutable" yaml:"serverexecutable"`
}
