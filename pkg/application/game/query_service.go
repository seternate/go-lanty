package game

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/seternate/go-lanty/pkg/application/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
	"github.com/seternate/go-lanty/pkg/domain/game"
)

type queryRow struct {
	GameSlug      string       `db:"game_slug"`
	GameName      string       `db:"game_name"`
	GameCreatedAt sql.NullTime `db:"game_created_at"`

	ExecID            sql.NullString `db:"exec_id"`
	ExecGameSlug      sql.NullString `db:"exec_game_slug"`
	ExecRole          sql.NullString `db:"exec_role"`
	ExecPath          sql.NullString `db:"exec_path"`
	ExecRequiresAdmin sql.NullBool   `db:"exec_requires_admin"`
	ExecFormat        sql.NullString `db:"exec_format"`
	ExecArgSeperator  sql.NullString `db:"exec_arg_seperator"`
	ExecCreatedAt     sql.NullTime   `db:"exec_created_at"`

	ArgID             sql.NullString  `db:"arg_id"`
	ArgGameExecID     sql.NullString  `db:"arg_game_exec_id"`
	ArgRole           sql.NullString  `db:"arg_role"`
	ArgName           sql.NullString  `db:"arg_name"`
	ArgRequired       sql.NullBool    `db:"arg_required"`
	ArgEnabled        sql.NullBool    `db:"arg_enabled"`
	ArgFormat         sql.NullString  `db:"arg_format"`
	ArgSeparator      sql.NullString  `db:"arg_separator"`
	ArgArg            sql.NullString  `db:"arg_arg"`
	ArgDescription    sql.NullString  `db:"arg_description"`
	ArgDefaultString  sql.NullString  `db:"arg_default_string"`
	ArgDefaultBool    sql.NullBool    `db:"arg_default_bool"`
	ArgDefaultInt     sql.NullInt64   `db:"arg_default_int"`
	ArgDefaultFloat   sql.NullFloat64 `db:"arg_default_float"`
	ArgEnums          pq.StringArray  `db:"arg_enums"`
	ArgMinInt         sql.NullInt64   `db:"arg_min_int"`
	ArgMaxInt         sql.NullInt64   `db:"arg_max_int"`
	ArgMinFloat       sql.NullFloat64 `db:"arg_min_float"`
	ArgMaxFloat       sql.NullFloat64 `db:"arg_max_float"`
	ArgFloatPrecision sql.NullInt64   `db:"arg_float_precision"`
	ArgOrderIndex     sql.NullInt64   `db:"arg_order_index"`
	ArgCreatedAt      sql.NullTime    `db:"arg_created_at"`
}

type queryServiceImpl struct {
	db               *sqlx.DB
	iconAssetService *asset.Service
	blobAssetService *asset.Service
}

func NewQueryService(db *sqlx.DB, iconAssetService *asset.Service, blobAssetService *asset.Service) QueryService {
	return &queryServiceImpl{
		db:               db,
		iconAssetService: iconAssetService,
		blobAssetService: blobAssetService,
	}
}

func (service *queryServiceImpl) GetGames() (GamesView, error) {
	query, args, err := sq.
		Select(
			"g.slug AS game_slug",
			"g.name AS game_name",
			"g.created_at AS game_created_at",
			"e.id AS exec_id",
			"e.game_slug AS exec_game_slug",
			"e.role AS exec_role",
			"e.path AS exec_path",
			"e.requires_admin AS exec_requires_admin",
			"e.format AS exec_format",
			"e.arg_seperator AS exec_arg_seperator",
			"e.created_at AS exec_created_at",
			"a.id AS arg_id",
			"a.game_exec_id AS arg_game_exec_id",
			"a.role AS arg_role",
			"a.name AS arg_name",
			"a.required AS arg_required",
			"a.enabled AS arg_enabled",
			"a.format AS arg_format",
			"a.arg_separator AS arg_separator",
			"a.arg AS arg_arg",
			"a.description AS arg_description",
			"a.default_string AS arg_default_string",
			"a.default_bool AS arg_default_bool",
			"a.default_int AS arg_default_int",
			"a.default_float AS arg_default_float",
			"a.enums AS arg_enums",
			"a.min_int AS arg_min_int",
			"a.max_int AS arg_max_int",
			"a.min_float AS arg_min_float",
			"a.max_float AS arg_max_float",
			"a.float_precision AS arg_float_precision",
			"a.order_index AS arg_order_index",
			"a.created_at AS arg_created_at",
		).
		From("games g").
		LeftJoin("game_execs e ON g.slug = e.game_slug").
		LeftJoin("game_args a ON e.id = a.game_exec_id").
		OrderBy("g.slug", "e.role", "a.order_index").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for games table: %w", err)
	}

	var rows []queryRow
	err = service.db.Select(&rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for games table: %s (%v): %w", query, args, err)
	}

	return assembleGamesView(rows), nil
}

func (service *queryServiceImpl) GetGame(slug string) (*GameView, error) {
	query, args, err := sq.
		Select(
			"g.slug AS game_slug",
			"g.name AS game_name",
			"g.created_at AS game_created_at",
			"e.game_slug AS exec_game_slug",
			"e.role AS exec_role",
			"e.path AS exec_path",
			"e.requires_admin AS exec_requires_admin",
			"e.format AS exec_format",
			"e.arg_seperator AS exec_arg_seperator",
			"e.created_at AS exec_created_at",
			"a.id AS arg_id",
			"a.game_exec_id AS arg_game_exec_id",
			"a.role AS arg_role",
			"a.name AS arg_name",
			"a.required AS arg_required",
			"a.enabled AS arg_enabled",
			"a.format AS arg_format",
			"a.arg_separator AS arg_separator",
			"a.arg AS arg_arg",
			"a.description AS arg_description",
			"a.default_string AS arg_default_string",
			"a.default_bool AS arg_default_bool",
			"a.default_int AS arg_default_int",
			"a.default_float AS arg_default_float",
			"a.enums AS arg_enums",
			"a.min_int AS arg_min_int",
			"a.max_int AS arg_max_int",
			"a.min_float AS arg_min_float",
			"a.max_float AS arg_max_float",
			"a.float_precision AS arg_float_precision",
			"a.order_index AS arg_order_index",
			"a.created_at AS arg_created_at",
		).
		From("games g").
		LeftJoin("game_execs e ON g.slug = e.game_slug").
		LeftJoin("game_args a ON e.id = a.game_exec_id").
		Where(sq.Eq{"g.slug": slug}).
		OrderBy("e.role", "a.order_index").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for games table for slug=%s: %w", slug, err)
	}

	var rows []queryRow
	err = service.db.Select(&rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for games table for slug=%s: %s (%v): %w", slug, query, args, err)
	}

	games := assembleGamesView(rows)
	if len(games) == 0 {
		return nil, domainerr.NotFoundErr("game", slug)
	} else if len(games) > 1 {
		return nil, fmt.Errorf("multiple games found for slug=%s", slug)
	}

	return &games[0], nil
}

func assembleGamesView(rows []queryRow) GamesView {
	if len(rows) == 0 {
		return GamesView{}
	}

	gamesMap := make(map[string]*GameView)

	for _, r := range rows {
		game, exists := gamesMap[r.GameSlug]
		if !exists {
			var createdAt time.Time
			if r.GameCreatedAt.Valid {
				createdAt = r.GameCreatedAt.Time
			}
			game = &GameView{
				Slug:      r.GameSlug,
				Name:      r.GameName,
				Execs:     make(map[string]GameExecView),
				CreatedAt: createdAt,
			}
			gamesMap[r.GameSlug] = game
		}

		if r.ExecID.Valid {
			execRole := r.ExecRole.String
			exec, execExists := game.Execs[execRole]
			if !execExists {
				var execCreatedAt time.Time
				if r.ExecCreatedAt.Valid {
					execCreatedAt = r.ExecCreatedAt.Time
				}
				exec = GameExecView{
					GameSlug:  r.ExecGameSlug.String,
					Role:      execRole,
					Path:      r.ExecPath.String,
					Args:      []GameArgView{},
					CreatedAt: execCreatedAt,
				}

				if r.ExecRequiresAdmin.Valid {
					exec.RequiresAdmin = &r.ExecRequiresAdmin.Bool
				}
				if r.ExecFormat.Valid {
					exec.Format = &r.ExecFormat.String
				}
				if r.ExecArgSeperator.Valid {
					exec.ArgSeperator = &r.ExecArgSeperator.String
				}

				game.Execs[execRole] = exec
			}

			if r.ArgID.Valid {
				var argCreatedAt time.Time
				if r.ArgCreatedAt.Valid {
					argCreatedAt = r.ArgCreatedAt.Time
				}
				arg := GameArgView{
					Role:       r.ArgRole.String,
					Name:       r.ArgName.String,
					Arg:        r.ArgArg.String,
					OrderIndex: r.ArgOrderIndex.Int64,
					CreatedAt:  argCreatedAt,
				}

				if r.ArgRequired.Valid {
					arg.Required = &r.ArgRequired.Bool
				}
				if r.ArgEnabled.Valid {
					arg.Enabled = &r.ArgEnabled.Bool
				}
				if r.ArgFormat.Valid {
					arg.Format = &r.ArgFormat.String
				}
				if r.ArgSeparator.Valid {
					arg.ArgSeparator = &r.ArgSeparator.String
				}
				if r.ArgDescription.Valid {
					arg.Description = &r.ArgDescription.String
				}
				if r.ArgDefaultString.Valid {
					arg.DefaultString = &r.ArgDefaultString.String
				}
				if r.ArgDefaultBool.Valid {
					arg.DefaultBool = &r.ArgDefaultBool.Bool
				}
				if r.ArgDefaultInt.Valid {
					arg.DefaultInt = &r.ArgDefaultInt.Int64
				}
				if r.ArgDefaultFloat.Valid {
					arg.DefaultFloat = &r.ArgDefaultFloat.Float64
				}
				if len(r.ArgEnums) > 0 {
					arg.Enums = []string(r.ArgEnums)
				}
				if r.ArgMinInt.Valid {
					arg.MinInt = &r.ArgMinInt.Int64
				}
				if r.ArgMaxInt.Valid {
					arg.MaxInt = &r.ArgMaxInt.Int64
				}
				if r.ArgMinFloat.Valid {
					arg.MinFloat = &r.ArgMinFloat.Float64
				}
				if r.ArgMaxFloat.Valid {
					arg.MaxFloat = &r.ArgMaxFloat.Float64
				}
				if r.ArgFloatPrecision.Valid {
					arg.FloatPrecision = &r.ArgFloatPrecision.Int64
				}

				exec.Args = append(exec.Args, arg)
				game.Execs[execRole] = exec
			}
		}
	}

	result := make(GamesView, 0, len(gamesMap))
	for _, game := range gamesMap {
		result = append(result, *game)
	}

	return result
}

func (service *queryServiceImpl) FetchIcon(slug string) (*AssetContent, error) {
	query, args, err := sq.
		Select("asset_id").
		From("game_assets").
		Where(sq.Eq{"game_slug": slug, "role": game.GAME_ASSET_ROLE_ICON.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for game assets table for slug=%s: %w", slug, err)
	}

	var assetID uuid.UUID
	err = service.db.Get(&assetID, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("game icon", slug).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for game assets table for slug=%s: %s (%v): %w", slug, query, args, err)
	}

	assetView, err := service.iconAssetService.Query.GetAsset(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get icon asset: %w", err)
	}

	data, err := service.iconAssetService.Query.GetAssetContent(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get icon asset content: %w", err)
	}

	return &AssetContent{
		Data:      data,
		Size:      assetView.Size,
		Checksum:  assetView.Checksum,
		Algorithm: assetView.Algorithm,
		MimeType:  assetView.MimeType,
	}, nil
}

func (service *queryServiceImpl) FetchBlob(slug string) (*AssetContent, error) {
	query, args, err := sq.
		Select("asset_id").
		From("game_assets").
		Where(sq.Eq{"game_slug": slug, "role": game.GAME_ASSET_ROLE_BLOB.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for game assets table for slug=%s: %w", slug, err)
	}

	var assetID uuid.UUID
	err = service.db.Get(&assetID, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("game blob", slug).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for game assets table for slug=%s: %s (%v): %w", slug, query, args, err)
	}

	assetView, err := service.blobAssetService.Query.GetAsset(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob asset: %w", err)
	}

	data, err := service.blobAssetService.Query.GetAssetContent(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob asset content: %w", err)
	}

	return &AssetContent{
		Data:      data,
		Size:      assetView.Size,
		Checksum:  assetView.Checksum,
		Algorithm: assetView.Algorithm,
		MimeType:  assetView.MimeType,
	}, nil
}
