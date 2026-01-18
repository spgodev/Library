package repository

import (
	"context"
	"database/sql"
	"time"

	"library/internal/domain"
)

type ReadingRepository struct {
	db *sql.DB
}

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{db: db} //2
}

func (r *ReadingRepository) Exists(ctx context.Context, bookID int, userID int64) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM book_readings
			WHERE book_id = $1 AND user_id = $2
		)
	`, bookID, userID).Scan(&ok)
	return ok, err
}

func (r *ReadingRepository) Insert(ctx context.Context, bookID int, userID int64, date time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO book_readings (book_id, user_id, read_date)
		VALUES ($1,$2,$3::date)
	`, bookID, userID, date)
	return err
}

func (r *ReadingRepository) ListByBook(ctx context.Context, bookID int) ([]domain.ReadingInfo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.full_name, br.read_date
		FROM book_readings br
		JOIN users u ON u.id = br.user_id
		WHERE br.book_id = $1
		ORDER BY br.read_date ASC, u.full_name ASC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.ReadingInfo
	for rows.Next() {
		var ri domain.ReadingInfo
		if err := rows.Scan(&ri.User, &ri.Date); err != nil {
			return nil, err
		}
		out = append(out, ri)
	}
	return out, rows.Err()
}
