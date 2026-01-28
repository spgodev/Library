package repository

import (
	"context"
	"database/sql"
	"errors"
	"library/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetIDByName(ctx context.Context, fullName string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, full_name FROM users WHERE full_name = $1
	`, fullName).Scan(&u.ID, &u.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) Insert(ctx context.Context, fullName string) (domain.User, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (full_name)
		VALUES ($1)
		RETURNING id
		    `, fullName).Scan(&id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{ID: id, Name: fullName}, nil
}
