package asset

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/url"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

var _ QueryService = (*queryServiceImpl)(nil)

type queryServiceImpl struct {
	db      *sqlx.DB
	storage StorageAdapter
}

func NewQueryService(db *sqlx.DB, storage StorageAdapter) QueryService {
	return &queryServiceImpl{
		db:      db,
		storage: storage,
	}
}

func (service *queryServiceImpl) GetAsset(id uuid.UUID) (*AssetView, error) {
	query, args, err := sq.
		Select("id", "url", "size", "checksum", "algorithm", "mime_type", "created_at").
		From("assets").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for assets table for asset id=%s: %w", id.String(), err)
	}

	view := &AssetView{}
	err = service.db.Get(view, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("asset", id.String()).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for assets table for asset id=%s: %s (%v): %w", id.String(), query, args, err)
	}

	return view, nil
}

func (service *queryServiceImpl) GetAssetContent(id uuid.UUID) (data io.ReadCloser, err error) {
	query, args, err := sq.
		Select("url").
		From("assets").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for assets table for asset id=%s: %w", id.String(), err)
	}

	u := ""
	err = service.db.Get(&u, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("asset", id.String()).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for assets table for asset id=%s: %s (%v): %w", id.String(), query, args, err)
	}

	assetURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse asset URL=%s: %w", u, err)
	}

	data, _, err = service.storage.Fetch(*assetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset from storage for asset id=%s: %w", id.String(), err)
	}

	return data, nil
}
