package game

import (
	"io"
	"time"

	"github.com/seternate/go-lanty/internal/domain/game"
)

type CommandService interface {
	UpsertGame(cmd UpsertGameCommand) (*game.Game, bool, error)
	DeleteGame(slug string) error
	StoreNewAsset(cmd StoreNewAssetCommand) (bool, error)
}

type UpsertGameCommand struct {
	Slug        string
	Name        string
	Executables []UpsertGameExecutable
}

type UpsertGameExecutable struct {
	Role              string
	Path              string
	RequiresAdmin     *bool
	Format            *string
	ArgumentSeperator *string
	Args              []UpsertGameExecutableArg
}

type UpsertGameExecutableArg struct {
	Role              string
	Name              string
	Required          *bool
	Enabled           *bool
	Format            *string
	ArgumentSeparator *string
	Argument          string
	Description       *string
	DefaultString     *string
	DefaultBool       *bool
	DefaultInt        *int64
	DefaultFloat      *float64
	EnumValues        []string
	MinInt            *int64
	MaxInt            *int64
	MinFloat          *float64
	MaxFloat          *float64
	FloatPrecision    *int64
	OrderIndex        int64
}

type StoreNewAssetCommand struct {
	Slug      string
	Role      game.GameAssetRole
	Checksum  string
	Algorithm string
	Data      io.Reader
}

type QueryService interface {
	GetGames() (GamesView, error)
	GetGame(slug string) (*GameView, error)
	FetchIcon(slug string) (*AssetContent, error)
	FetchBlob(slug string) (*AssetContent, error)
}

type GameView struct {
	Slug      string
	Name      string
	CreatedAt time.Time
	Execs     map[string]GameExecView
}

type GamesView []GameView

type GameExecView struct {
	GameSlug      string
	Role          string
	Path          string
	RequiresAdmin *bool
	Format        *string
	ArgSeperator  *string
	CreatedAt     time.Time
	Args          []GameArgView
}

type GameArgView struct {
	Role           string
	Name           string
	Required       *bool
	Enabled        *bool
	Format         *string
	ArgSeparator   *string
	Arg            string
	Description    *string
	DefaultString  *string
	DefaultBool    *bool
	DefaultInt     *int64
	DefaultFloat   *float64
	Enums          []string
	MinInt         *int64
	MaxInt         *int64
	MinFloat       *float64
	MaxFloat       *float64
	FloatPrecision *int64
	OrderIndex     int64
	CreatedAt      time.Time
}

type AssetContent struct {
	Size      uint64
	Checksum  string
	Algorithm string
	MimeType  string
	Data      io.ReadCloser
}

type Service struct {
	Query   QueryService
	Command CommandService
}

func NewService(query QueryService, command CommandService) *Service {
	return &Service{
		Query:   query,
		Command: command,
	}
}
