package story

import (
	"context"
	"fmt"
	"time"

	"library/internal/domain"
)

// интерфейсы вместо конкретных struct-репозиториев (для тестов)
type BookRepo interface {
	ExistsByID(ctx context.Context, id int) (bool, error)
	Insert(ctx context.Context, b domain.Book) error
	ExistsByTitle(ctx context.Context, title string) (bool, error)
	FindByTitle(ctx context.Context, title string) ([]domain.Book, error)
	ListByYear(ctx context.Context, year int) ([]domain.Book, error)
	ListByAuthor(ctx context.Context, author string) ([]domain.Book, error)
	ListSortedByYear(ctx context.Context, asc bool) ([]domain.Book, error)
}

type UserRepo interface {
	GetIDByName(ctx context.Context, fullName string) (int64, bool, error)
	Insert(ctx context.Context, fullName string) (int64, error)
}

type ReadingRepo interface {
	Exists(ctx context.Context, bookID int, userID int64) (bool, error)
	Insert(ctx context.Context, bookID int, userID int64, date time.Time) error
	ListByBook(ctx context.Context, bookID int) ([]domain.ReadingInfo, error)
}

type LibraryStory struct {
	books    BookRepo
	users    UserRepo
	readings ReadingRepo
}

func NewLibraryStory(books BookRepo, users UserRepo, readings ReadingRepo) *LibraryStory {
	return &LibraryStory{books: books, users: users, readings: readings}
}

func (s *LibraryStory) AddBook(ctx context.Context, book domain.Book) error {
	exists, err := s.books.ExistsByID(ctx, book.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("book with ID %d already exists", book.ID)
	}
	return s.books.Insert(ctx, book)
}

func (s *LibraryStory) HasBookTitle(ctx context.Context, title string) (bool, error) {
	return s.books.ExistsByTitle(ctx, title)
}

func (s *LibraryStory) FindBookByTitle(ctx context.Context, title string) ([]domain.Book, error) {
	return s.books.FindByTitle(ctx, title)
}

func (s *LibraryStory) MarkAsRead(ctx context.Context, bookID int, userFullName string, date time.Time) error {
	bookExists, err := s.books.ExistsByID(ctx, bookID)
	if err != nil {
		return err
	}
	if !bookExists {
		return fmt.Errorf("book not found")
	}

	uid, found, err := s.users.GetIDByName(ctx, userFullName)
	if err != nil {
		return err
	}
	if !found {
		uid, err = s.users.Insert(ctx, userFullName)
		if err != nil {
			return err
		}
	}

	already, err := s.readings.Exists(ctx, bookID, uid)
	if err != nil {
		return err
	}
	if already {
		return fmt.Errorf("user already marked as read")
	}

	return s.readings.Insert(ctx, bookID, uid, date)
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
	bookExists, err := s.books.ExistsByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if !bookExists {
		return nil, fmt.Errorf("book not found")
	}
	return s.readings.ListByBook(ctx, bookID)
}
