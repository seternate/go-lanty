package user

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/domain/user"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
)

var _ user.UserRepository = (*userrepository)(nil)

type userrepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.UserRepository {
	return &userrepository{
		db: db,
	}
}

func (repository *userrepository) GetUser(ipv4Address string) (*user.User, error) {
	query, args, err := sq.
		Select("ipv4_address", "username").
		From("users").
		Where(sq.Eq{"ipv4_address": ipv4Address}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for users table for user ipv4_address=%s: %w", ipv4Address, err)
	}

	userrow := &UserRow{}
	err = repository.db.Get(userrow, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerr.NotFoundErr("user", ipv4Address).WithCause(err)
		}
		return nil, fmt.Errorf("database query failed for users table for user ipv4_address=%s: %s (%v): %w", ipv4Address, query, args, err)
	}

	user, err := userrow.Assemble()
	if err != nil {
		return nil, fmt.Errorf("failed to assemble user: %w", err)
	}

	return user, nil
}

func (repository *userrepository) SaveUser(u *user.User) error {
	disassembledUser := Disassemble(u)

	columnvalues, err := database.GenerateColumnValueMap(*disassembledUser)
	if err != nil {
		return fmt.Errorf("failed to generate column value map for user ipv4_address=%s: %w", u.IPv4Address.String(), err)
	}

	// Use UPSERT (INSERT ... ON CONFLICT ... DO UPDATE) for SaveUser
	query, args, err := sq.
		Insert("users").
		SetMap(columnvalues).
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (ipv4_address) DO UPDATE SET username = EXCLUDED.username").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for users table for user ipv4_address=%s: %w", u.IPv4Address.String(), err)
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		err = fmt.Errorf("database query failed for users table for user ipv4_address=%s: %s (%v): %w", u.IPv4Address.String(), query, args, err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domainerr.ConflictErr("user", u.IPv4Address.String()).WithCause(err)
			}
		}
		return err
	}

	return nil
}

func (repository *userrepository) DeleteUser(ipv4Address string) error {
	query, args, err := sq.
		Delete("users").
		Where(sq.Eq{"ipv4_address": ipv4Address}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for users table for user ipv4_address=%s: %w", ipv4Address, err)
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed for users table for user ipv4_address=%s: %s (%v): %w", ipv4Address, query, args, err)
	}

	return nil
}
