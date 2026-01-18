package repository

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetIDByName(ctx context.Context, fullName string) (int64, bool, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM users WHERE full_name = $1
	`, fullName).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func (r *UserRepository) Insert(ctx context.Context, fullName string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (full_name)
		VALUES ($1)
		RETURNING id
	`, fullName).Scan(&id)
	return id, err
}
