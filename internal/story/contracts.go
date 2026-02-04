package story

import (
	"context"
	"time"

	"library/internal/domain"
)

type BookRepo interface {
	GetAll(ctx context.Context) ([]domain.Book, error)
	Insert(ctx context.Context, b domain.Book) (domain.Book, error)

	FindAllByTitle(ctx context.Context, title string) ([]domain.Book, error)
	FindAllByYear(ctx context.Context, year int) ([]domain.Book, error)
	FindAllByAuthor(ctx context.Context, author string) ([]domain.Book, error)

	ListSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error)
}

type UserRepo interface {
	GetUserByName(ctx context.Context, fullName string) (domain.User, error)
	Insert(ctx context.Context, fullName string) (domain.User, error)
}

type ReadingRepo interface {
	Insert(ctx context.Context, bookID int64, userID int64, date time.Time) (domain.BookReading, error)
	ListByBook(ctx context.Context, bookID int64) ([]domain.ReadingInfo, error)
}
