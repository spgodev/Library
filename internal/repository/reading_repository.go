package repository

import (
	"context"
	"time"

	"library/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadingRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ReadingRepository {
	return &ReadingRepository{db: db}
}

func (r *ReadingRepository) Insert(ctx context.Context, bookID int64, userID int64, date time.Time) (domain.BookReading, error) {
	var br domain.BookReading

	err := r.db.QueryRow(ctx, `
		INSERT INTO book_readings (book_id, user_id, read_date)
		VALUES ($1,$2,$3::date)
		RETURNING book_id, user_id, read_date
	`, bookID, userID, date).Scan(&br.BookID, &br.UserID, &br.ReadDate)

	if err != nil {
		return domain.BookReading{}, err
	}
	return br, nil
}

func (r *ReadingRepository) ListByBook(ctx context.Context, bookID int64) ([]domain.ReadingInfo, error) {
	rows, err := r.db.Query(ctx, `
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

	out := make([]domain.ReadingInfo, 0)
	for rows.Next() {
		var ri domain.ReadingInfo
		if err := rows.Scan(&ri.User, &ri.Date); err != nil {
			return nil, err
		}
		out = append(out, ri)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
