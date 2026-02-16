package repository

import (
	"context"
	"fmt"
	"strings"

	"library/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *BookRepository {
	return &BookRepository{db: db}
}

func scanBooks(rows pgx.Rows) ([]domain.Book, error) {
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.LibraryID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *BookRepository) Insert(ctx context.Context, b domain.Book) (domain.Book, error) {
	var out domain.Book
	err := r.db.QueryRow(ctx, `
		INSERT INTO books (library_id, title, author, year, pages)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, library_id, title, author, year, pages
	`, b.LibraryID, b.Title, b.Author, b.Year, b.Pages).Scan(
		&out.ID, &out.LibraryID, &out.Title, &out.Author, &out.Year, &out.Pages,
	)
	if err != nil {
		return domain.Book{}, err
	}
	return out, nil
}

func (r *BookRepository) FindAllByTitle(ctx context.Context, title string) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE title = $1
		ORDER BY id
	`, title)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) FindAllByYear(ctx context.Context, year int) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE year = $1
		ORDER BY id
	`, year)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) FindAllByAuthor(ctx context.Context, author string) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE author = $1
		ORDER BY id
	`, author)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) ListSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error) {
	order := "ASC"
	if !asc {
		order = "DESC"
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		ORDER BY year `+order+`, id ASC
	`)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) GetAllByLibraryID(ctx context.Context, libraryID int64) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE library_id = $1
		ORDER BY id
	`, libraryID)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) FindAllByAuthorInLibrary(ctx context.Context, libraryID int64, author string) ([]domain.Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE library_id = $1 AND author = $2
		ORDER BY id
	`, libraryID, author)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}
func (r *BookRepository) FindByFilters(ctx context.Context, libraryID int64, year *int, authorSubstr *string, titleSubstr *string) ([]domain.Book, error) {
	query := `
		SELECT id, library_id, title, author, year, pages
		FROM books
		WHERE library_id = $1
	`
	args := []any{libraryID}
	n := 1

	if year != nil {
		n++
		query += fmt.Sprintf(" AND year = $%d", n)
		args = append(args, *year)
	}
	if authorSubstr != nil && strings.TrimSpace(*authorSubstr) != "" {
		n++
		query += fmt.Sprintf(" AND author ILIKE $%d", n)
		args = append(args, "%"+strings.TrimSpace(*authorSubstr)+"%")
	}
	if titleSubstr != nil && strings.TrimSpace(*titleSubstr) != "" {
		n++
		query += fmt.Sprintf(" AND title ILIKE $%d", n)
		args = append(args, "%"+strings.TrimSpace(*titleSubstr)+"%")
	}

	query += " ORDER BY id"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}

func (r *BookRepository) ExistsInLibrary(ctx context.Context, libraryID, bookID int64) (bool, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		SELECT id
		FROM books
		WHERE id = $1 AND library_id = $2
	`, bookID, libraryID).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
