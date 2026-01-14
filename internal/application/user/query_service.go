package user

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ QueryService = (*queryServiceImpl)(nil)

type queryRow struct {
	IPv4Address string       `db:"ipv4_address"`
	Username    string       `db:"username"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

type queryServiceImpl struct {
	db *sqlx.DB
}

func NewQueryService(db *sqlx.DB) QueryService {
	return &queryServiceImpl{
		db: db,
	}
}

func (service *queryServiceImpl) GetUsers() (UsersView, error) {
	query, args, err := sq.
		Select("ipv4_address", "username", "created_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for users table: %w", err)
	}

	var rows []queryRow
	err = service.db.Select(&rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed for users table: %s (%v): %w", query, args, err)
	}

	views := make(UsersView, 0, len(rows))
	for _, row := range rows {
		view := UserView{
			IPv4Address: row.IPv4Address,
			Username:    row.Username,
		}
		if row.CreatedAt.Valid {
			view.CreatedAt = row.CreatedAt.Time
		}
		views = append(views, view)
	}

	return views, nil
}
