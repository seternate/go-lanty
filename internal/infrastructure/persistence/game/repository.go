package game

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/domain/game"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
)

var _ game.GameRepository = (*gamerepository)(nil)

type gamerepository struct {
	db *sqlx.DB
}

func NewGameRepository(db *sqlx.DB) game.GameRepository {
	return &gamerepository{
		db: db,
	}
}

func (repository *gamerepository) GetGame(slug string) (*game.Game, error) {
	gamerow, err := repository.fetchGameRow(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game row: %w", err)
	}

	execrows, err := repository.fetchExecRows(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exec rows: %w", err)
	}

	execids := make([]uuid.UUID, len(execrows))
	for _, execrow := range execrows {
		execids = append(execids, execrow.ID)
	}

	argrows, err := repository.fetchArgRows(execids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch arg rows for game slug=%s: %w", slug, err)
	}

	assetrows, err := repository.fetchAssetRows(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset rows: %w", err)
	}

	game, err := gamerow.Assemble(execrows, argrows, assetrows)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble game: %w", err)
	}

	return game, nil
}

func (repository *gamerepository) SaveGame(g *game.Game) error {
	disassembledGame := Disassemble(g)

	tx, err := repository.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for game slug=%s: %w", g.Slug, err)
	}

	err = deleteGameRow(tx, g.Slug)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete game: %w", err)
	}

	err = saveGameRow(tx, disassembledGame.GameRow)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save game for slug=%s: %w", g.Slug, err)
	}

	err = saveExecRows(tx, disassembledGame.ExecRows)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save game exec rows for game slug=%s: %w", g.Slug, err)
	}

	err = saveArgRows(tx, disassembledGame.ArgRows)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save game arg rows for game slug=%s: %w", g.Slug, err)
	}

	err = saveAssetRows(tx, disassembledGame.AssetRows)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save game asset rows for game slug=%s: %w", g.Slug, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for game slug=%s: %w", g.Slug, err)
	}

	return nil
}

func (repository *gamerepository) DeleteGame(slug string) error {
	err := deleteGameRow(repository.db, slug)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

func (repository *gamerepository) fetchGameRow(slug string) (*GameRow, error) {
	query, args, err := sq.
		Select("slug", "name").
		From("games").
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for games table for game slug=%s: %w", slug, err)
	}

	gamerow := &GameRow{}
	err = repository.db.Get(gamerow, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("game", slug).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for games table for game slug=%s: %s (%v): %w", slug, query, args, err)
	}

	return gamerow, nil
}

func (repository *gamerepository) fetchExecRows(slug string) ([]GameExecRow, error) {
	query, args, err := sq.
		Select("id", "game_slug", "role", "path", "requires_admin", "format", "arg_seperator").
		From("game_execs").
		Where(sq.Eq{"game_slug": slug}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for game_execs table for game slug=%s: %w", slug, err)
	}

	execrows := make([]GameExecRow, 0)
	err = repository.db.Select(&execrows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for game_execs table for game slug=%s: %s (%v): %w", slug, query, args, err)
	}

	return execrows, nil
}

func (repository *gamerepository) fetchArgRows(execids []uuid.UUID) ([]GameArgRow, error) {
	if len(execids) == 0 {
		return nil, nil
	}

	query, args, err := sq.
		Select(
			"game_exec_id",
			"role",
			"name",
			"required",
			"enabled",
			"format",
			"arg_separator",
			"arg",
			"description",
			"default_string",
			"default_bool",
			"default_int",
			"default_float",
			"enums",
			"min_int",
			"max_int",
			"min_float",
			"max_float",
			"float_precision",
			"order_index",
		).
		From("game_args").
		Where(sq.Eq{"game_exec_id": execids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for game_args table for game exec ids=%v: %w", execids, err)
	}

	argrows := make([]GameArgRow, 0)
	err = repository.db.Select(&argrows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for game_args table for game exec ids=%v: %s (%v): %w", execids, query, args, err)
	}

	return argrows, nil
}

func (repository *gamerepository) fetchAssetRows(slug string) ([]GameAssetRow, error) {
	query, args, err := sq.
		Select("game_slug", "asset_id", "role").
		From("game_assets").
		Where(sq.Eq{"game_slug": slug}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for game_assets table for game slug=%s: %w", slug, err)
	}

	assetrows := make([]GameAssetRow, 0)
	err = repository.db.Select(&assetrows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for game_assets table for game slug=%s: %s (%v): %w", slug, query, args, err)
	}

	return assetrows, nil
}

func deleteGameRow(db sq.BaseRunner, slug string) error {
	query, args, err := sq.
		Delete("games").
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for games table for game slug=%s: %w", slug, err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed for games table for game slug=%s: %s (%v): %w", slug, query, args, err)
	}

	return nil
}

func saveRows[T any](db sq.BaseRunner, tableName string, rows []T) error {
	if len(rows) == 0 {
		return nil
	}

	sample := rows[0]
	cols, err := database.GenerateColumnValueMap(sample)
	if err != nil {
		return fmt.Errorf("failed to generate column value map for sample %T in %s table: %w", sample, tableName, err)
	}

	colNames := make([]string, 0, len(cols))
	for col := range cols {
		colNames = append(colNames, col)
	}
	sort.Strings(colNames)

	qb := sq.
		Insert(tableName).
		Columns(colNames...)

	for _, row := range rows {
		m, err := database.GenerateColumnValueMap(row)
		if err != nil {
			return fmt.Errorf("failed to generate column value map for %s table: %w", tableName, err)
		}
		vals := make([]interface{}, 0, len(colNames))
		for _, col := range colNames {
			vals = append(vals, m[col])
		}
		qb = qb.Values(vals...)
	}

	query, args, err := qb.
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for %s table: %w", tableName, err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed for %s table: %s (%v): %w", tableName, query, args, err)
	}

	return nil
}

func saveGameRow(db sq.BaseRunner, gamerow GameRow) error {
	columnvalues, err := database.GenerateColumnValueMap(gamerow)
	if err != nil {
		return fmt.Errorf("failed to generate column value map for game row: %w", err)
	}

	query, args, err := sq.
		Insert("games").
		SetMap(columnvalues).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for games table for game slug=%s: %w", gamerow.Slug, err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domainerr.ConflictErr("game", gamerow.Slug).WithCause(err)
			}
		}
		return fmt.Errorf("database query failed for games table for game slug=%s: %s (%v): %w", gamerow.Slug, query, args, err)
	}

	return nil
}

func saveExecRows(db sq.BaseRunner, execrows []GameExecRow) error {
	if err := saveRows(db, "game_execs", execrows); err != nil {
		return fmt.Errorf("failed to save rows: %w", err)
	}
	return nil
}

func saveArgRows(db sq.BaseRunner, argrows []GameArgRow) error {
	if err := saveRows(db, "game_args", argrows); err != nil {
		return fmt.Errorf("failed to save rows: %w", err)
	}
	return nil
}

func saveAssetRows(db sq.BaseRunner, assetrows []GameAssetRow) error {
	if err := saveRows(db, "game_assets", assetrows); err != nil {
		return fmt.Errorf("failed to save rows: %w", err)
	}
	return nil
}
