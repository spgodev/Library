package story

import (
	"context"
	"errors"
	"time"

	"library/internal/domain"
)

type LibraryStory struct {
	books    BookRepo
	users    UserRepo
	readings ReadingRepo
}

func New(books BookRepo, users UserRepo, readings ReadingRepo) *LibraryStory {
	return &LibraryStory{books: books, users: users, readings: readings}
}

func (s *LibraryStory) AddBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	return s.books.Insert(ctx, book)
}

func (s *LibraryStory) HasBookTitle(ctx context.Context, title string) (bool, error) {
	books, err := s.books.FindAllByTitle(ctx, title)
	if err != nil {
		return false, err
	}
	return len(books) > 0, nil
}

func (s *LibraryStory) FindBookByTitle(ctx context.Context, title string) ([]domain.Book, error) {
	return s.books.FindAllByTitle(ctx, title)
}

func (s *LibraryStory) MarkAsRead(ctx context.Context, bookID int64, userFullName string, date time.Time) (domain.BookReading, error) {
	u, err := s.users.GetUserByName(ctx, userFullName)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			u, err = s.users.Insert(ctx, userFullName)
			if err != nil {
				return domain.BookReading{}, err
			}
		} else {
			return domain.BookReading{}, err
		}
	}

	return s.readings.Insert(ctx, bookID, u.ID, date)
}

func (s *LibraryStory) GetBooksByYear(ctx context.Context, year int) ([]domain.Book, error) {
	return s.books.FindAllByYear(ctx, year)
}

func (s *LibraryStory) GetBooksByAuthor(ctx context.Context, author string) ([]domain.Book, error) {
	return s.books.FindAllByAuthor(ctx, author)
}

func (s *LibraryStory) GetBooksSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error) {
	return s.books.ListSortedByYear(ctx, asc)
}

func (s *LibraryStory) GetReadersByBook(ctx context.Context, bookID int64) ([]domain.ReadingInfo, error) {
	return s.readings.ListByBook(ctx, bookID)
}

func (s *LibraryStory) LoadLibrary(ctx context.Context, name string) (*domain.Library, error) {
	books, err := s.books.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Library{Name: name, Books: books}, nil
}
