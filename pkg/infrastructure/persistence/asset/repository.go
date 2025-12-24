package asset

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/seternate/go-lanty/pkg/domain/asset"
	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
	"github.com/seternate/go-lanty/pkg/infrastructure/database"
)

var _ asset.AssetRepository = (*assetrepository)(nil)

type assetrepository struct {
	db *sqlx.DB
}

func NewAssetRepository(db *sqlx.DB) asset.AssetRepository {
	return &assetrepository{
		db: db,
	}
}

func (repository *assetrepository) GetAsset(id uuid.UUID) (*asset.Asset, error) {
	query, args, err := sq.
		Select("id", "url", "size", "checksum", "algorithm", "mime_type").
		From("assets").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for assets table for asset id=%s: %w", id.String(), err)
	}

	assetrow := &AssetRow{}
	err = repository.db.Get(assetrow, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.NotFoundErr("asset", id.String()).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for assets table for asset id=%s: %s (%v): %w", id.String(), query, args, err)
	}

	asset, err := assetrow.Assemble()
	if err != nil {
		return nil, fmt.Errorf("failed to assemble asset: %w", err)
	}

	return asset, nil
}

func (repository *assetrepository) CreateAsset(a *asset.Asset) error {
	disassambledAsset := Disassemble(a)

	columnvalues, err := database.GenerateColumnValueMap(*disassambledAsset)
	if err != nil {
		return fmt.Errorf("failed to generate column value map for asset id=%s: %w", a.ID.String(), err)
	}

	query, args, err := sq.
		Insert("assets").
		SetMap(columnvalues).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domainErrors.ConflictErr("asset", a.ID.String()).WithCause(err)
			}
		}
		return fmt.Errorf("failed to build query for assets table for asset id=%s: %w", a.ID.String(), err)
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed for assets table for asset id=%s: %s (%v): %w", a.ID.String(), query, args, err)
	}

	return nil
}

func (repository *assetrepository) DeleteAsset(id uuid.UUID) error {
	query, args, err := sq.
		Delete("assets").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for assets table for asset id=%s: %w", id.String(), err)
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed for assets table for asset id=%s: %s (%v): %w", id.String(), query, args, err)
	}

	return nil
}
