package story

import (
	"context"
	"errors"
	"time"

	"library/internal/domain"
)

type LibraryStory struct {
	books     BookRepo
	users     UserRepo
	readings  ReadingRepo
	libraries LibraryRepo
}

func New(books BookRepo, users UserRepo, readings ReadingRepo, libraries LibraryRepo) *LibraryStory {
	return &LibraryStory{books: books, users: users, readings: readings, libraries: libraries}
}

func (s *LibraryStory) CreateLibrary(ctx context.Context, name string) (domain.Library, error) {
	return s.libraries.Insert(ctx, name)
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

// 1
func (s *LibraryStory) LoadLibrary(ctx context.Context, name string) (*domain.Library, error) {
	books, err := s.books.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Library{Name: name, Books: books}, nil
}

// 2
func (s *LibraryStory) AddBookToLibrary(ctx context.Context, libraryID int64, b domain.Book) (domain.Book, error) {
	b.LibraryID = libraryID
	return s.books.Insert(ctx, b)
}

// 3
func (s *LibraryStory) AddUserToLibrary(ctx context.Context, libraryID int64, fullName string) (domain.User, error) {
	return s.users.InsertIntoLibrary(ctx, libraryID, fullName)
}

// 4
func (s *LibraryStory) GetBooksByLibrary(ctx context.Context, libraryID int64) ([]domain.Book, error) {
	return s.books.GetAllByLibraryID(ctx, libraryID)
}

// 5
func (s *LibraryStory) GetBooksByAuthorInLibrary(ctx context.Context, libraryID int64, author string) ([]domain.Book, error) {
	return s.books.FindAllByAuthorInLibrary(ctx, libraryID, author)
}

// 6
func (s *LibraryStory) GetBooksByFilters(ctx context.Context, libraryID int64, year *int, authorSubstr *string, titleSubstr *string) ([]domain.Book, error) {
	return s.books.FindByFilters(ctx, libraryID, year, authorSubstr, titleSubstr)
}

// 7
func (s *LibraryStory) AddReading(ctx context.Context, libraryID, userID, bookID int64, date time.Time) (domain.BookReading, error) {
	if _, err := s.users.GetByIDInLibrary(ctx, libraryID, userID); err != nil {
		return domain.BookReading{}, err
	}

	ok, err := s.books.ExistsInLibrary(ctx, libraryID, bookID)
	if err != nil {
		return domain.BookReading{}, err
	}
	if !ok {
		return domain.BookReading{}, domain.NotFoundError
	}

	return s.readings.Insert(ctx, bookID, userID, date)
}

// 8
func (s *LibraryStory) GetReadingsByUser(ctx context.Context, libraryID, userID int64) ([]domain.ReadBookItem, error) {
	_, err := s.users.GetByIDInLibrary(ctx, libraryID, userID)
	if err != nil {
		return nil, err
	}
	return s.readings.ListReadingsByUser(ctx, libraryID, userID)
}
