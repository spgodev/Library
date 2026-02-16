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

func (r *ReadingRepository) Insert(ctx context.Context, bookID, userID int64, date time.Time) (domain.BookReading, error) {
	var br domain.BookReading

	err := r.db.QueryRow(ctx, `
		INSERT INTO book_readings (book_id, user_id, read_date)
		VALUES ($1,$2,$3::date)
		ON CONFLICT (book_id, user_id, read_date) DO UPDATE SET read_date = EXCLUDED.read_date
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
func (r *ReadingRepository) GetReadBooksByUser(ctx context.Context, libraryID, userID int64) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT b.id, b.library_id, b.title, b.author, b.year, b.pages
		FROM book_readings br
		JOIN books b ON b.id = br.book_id
		WHERE br.user_id = $1 AND b.library_id = $2
		ORDER BY b.id
	`, userID, libraryID)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *ReadingRepository) ListReadingsByUser(ctx context.Context, libraryID, userID int64) ([]domain.ReadBookItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.library_id, b.title, b.author, b.year, b.pages, br.read_date
		FROM book_readings br
		JOIN books b ON b.id = br.book_id
		WHERE br.user_id = $1 AND b.library_id = $2
		ORDER BY br.read_date DESC, b.id ASC
	`, userID, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.ReadBookItem, 0)
	for rows.Next() {
		var item domain.ReadBookItem
		if err := rows.Scan(
			&item.Book.ID, &item.Book.LibraryID, &item.Book.Title, &item.Book.Author, &item.Book.Year, &item.Book.Pages,
			&item.ReadDate,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
