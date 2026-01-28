package story

import (
	"context"
	"errors"
	"time"

	"library/internal/domain"
)

// интерфейсы вместо конкретных struct-репозиториев (для тестов)
type BookRepo interface {
	ListAll(ctx context.Context) ([]domain.Book, error)
	ExistsByID(ctx context.Context, id int) (bool, error)
	Insert(ctx context.Context, b domain.Book) (domain.Book, error)
	ExistsByTitle(ctx context.Context, title string) (bool, error)
	FindByTitle(ctx context.Context, title string) ([]domain.Book, error)
	ListByYear(ctx context.Context, year int) ([]domain.Book, error)
	ListByAuthor(ctx context.Context, author string) ([]domain.Book, error)
	ListSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error)
}

type UserRepo interface {
	GetIDByName(ctx context.Context, fullName string) (domain.User, error)
	Insert(ctx context.Context, fullName string) (domain.User, error)
}

type ReadingRepo interface {
	Insert(ctx context.Context, bookID int, userID int64, date time.Time) error
	ListByBook(ctx context.Context, bookID int) ([]domain.ReadingInfo, error)
}

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
	books, err := s.books.FindByTitle(ctx, title)
	if err != nil {
		return false, err
	}
	return len(books) > 0, nil
}

func (s *LibraryStory) FindBookByTitle(ctx context.Context, title string) ([]domain.Book, error) {
	return s.books.FindByTitle(ctx, title)
}

func (s *LibraryStory) MarkAsRead(ctx context.Context, bookID int, userFullName string, date time.Time) error {
	u, err := s.users.GetIDByName(ctx, userFullName)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			u, err = s.users.Insert(ctx, userFullName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = s.readings.Insert(ctx, bookID, u.ID, date)
	return err
}

func (s *LibraryStory) GetBooksByYear(ctx context.Context, year int) ([]domain.Book, error) {
	return s.books.ListByYear(ctx, year)
}

func (s *LibraryStory) GetBooksByAuthor(ctx context.Context, author string) ([]domain.Book, error) {
	return s.books.ListByAuthor(ctx, author)
}

func (s *LibraryStory) GetBooksSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error) {
	return s.books.ListSortedByYear(ctx, asc)
}

func (s *LibraryStory) GetReadersByBook(ctx context.Context, bookID int) ([]domain.ReadingInfo, error) {
	return s.readings.ListByBook(ctx, bookID)
}

func (s *LibraryStory) LoadLibrary(ctx context.Context, name string) (*domain.Library, error) {
	books, err := s.books.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Library{
		Name:  name,
		Books: books,
	}, nil
}
