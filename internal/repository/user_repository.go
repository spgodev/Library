package repository

import (
	"context"
	"errors"

	"library/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByName(ctx context.Context, fullName string) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRow(ctx, `
		SELECT id, full_name
		FROM users
		WHERE full_name = $1
	`, fullName).Scan(&u.ID, &u.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) Insert(ctx context.Context, fullName string) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRow(ctx, `
		INSERT INTO users (full_name)
		VALUES ($1)
		RETURNING id, full_name
	`, fullName).Scan(&u.ID, &u.Name)

	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}
