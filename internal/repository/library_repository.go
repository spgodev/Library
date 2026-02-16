package repository

import (
	"context"

	"library/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LibraryRepository struct {
	db *pgxpool.Pool
}

func NewLibraryRepository(db *pgxpool.Pool) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) Insert(ctx context.Context, name string) (domain.Library, error) {
	var out domain.Library
	err := r.db.QueryRow(ctx, `
		INSERT INTO libraries (name)
		VALUES ($1)
		RETURNING id, name
	`, name).Scan(&out.ID, &out.Name)

	if err != nil {
		return domain.Library{}, err
	}
	return out, nil
}

func (r *UserRepository) InsertIntoLibrary(ctx context.Context, libraryID int64, fullName string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (full_name, library_id)
		VALUES ($1, $2)
		RETURNING id, full_name, library_id
	`, fullName, libraryID).Scan(&u.ID, &u.Name, &u.LibraryID)

	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
