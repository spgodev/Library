package repository

import (
	"context"
	"database/sql"

	"library/internal/domain"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)`, id,
	).Scan(&ok)
	return ok, err
}

func (r *BookRepository) ExistsByTitle(ctx context.Context, title string) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM books WHERE title = $1)
	`, title).Scan(&ok)
	return ok, err
}

func (r *BookRepository) Insert(ctx context.Context, b domain.Book) (domain.Book, error) {
	var out domain.Book
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO books (title, author, year, pages)
		VALUES ($1,$2,$3,$4)
		RETURNING id, title, author, year, pages
	`, b.Title, b.Author, b.Year, b.Pages).Scan(
		&out.ID, &out.Title, &out.Author, &out.Year, &out.Pages,
	)
	if err != nil {
		return domain.Book{}, err
	}
	return out, nil
}

func (r *BookRepository) FindByTitle(ctx context.Context, title string) ([]domain.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, year, pages
		FROM books
		WHERE title = $1
		ORDER BY id
	`, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BookRepository) ListByYear(ctx context.Context, year int) ([]domain.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, year, pages
		FROM books
		WHERE year = $1
		ORDER BY id
	`, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BookRepository) ListByAuthor(ctx context.Context, author string) ([]domain.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, year, pages
		FROM books
		WHERE author = $1
		ORDER BY id
	`, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BookRepository) ListSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error) {
	order := "ASC"
	if !asc {
		order = "DESC"
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, year, pages
		FROM books
		ORDER BY year `+order+`, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BookRepository) ListAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, year, pages
		FROM books
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Book, 0)
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Pages); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
