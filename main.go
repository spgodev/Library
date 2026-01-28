package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"library/internal/app"
	"library/internal/repository"
	"library/internal/story"
)

const defaultDSN = "postgres://postgres:postgres@localhost:5434/library?sslmode=disable"

func printBooks(title string, books any) {
	b, _ := json.MarshalIndent(books, "", "  ")
	fmt.Println("\n" + title + ":\n" + string(b))
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	conn, err := app.Open(ctx, "postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := app.RunMigrations(conn); err != nil {
		log.Fatal(err)
	}

	bookRepo := repository.NewBookRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	readingRepo := repository.New(conn)

	lib := story.New(bookRepo, userRepo, readingRepo)

	// 2) наличие по названию
	ok, err := lib.HasBookTitle(ctx, "1984")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Has '1984':", ok)

	// 2) поиск по названию
	found, err := lib.FindBookByTitle(ctx, "Design Patterns")
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Find by title = Design Patterns", found)

	// 3) кто читал книгу 1
	readers, err := lib.GetReadersByBook(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nReaders of book #1:")
	for _, r := range readers {
		fmt.Printf("  - %-12s %s\n", r.User, r.Date.Format("2006-01-02"))
	}

	// 4) отметить прочтение
	err = lib.MarkAsRead(ctx, 3, "Vasya Pupkin", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	fmt.Println("\nMarkAsRead(book=3, user=Vasya Pupkin):", err)

	// 5) книги за год
	byYear, err := lib.GetBooksByYear(ctx, 1949)
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Books by year = 1949", byYear)

	// 6) книги автора
	byAuthor, err := lib.GetBooksByAuthor(ctx, "George Orwell")
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Books by author = George Orwell", byAuthor)

	// 7) сортировка
	sorted, err := lib.GetBooksSortedByYear(ctx, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nSorted by year ASC:")
	fmt.Println("  ID | YEAR | TITLE")
	fmt.Println(" ----+------+----------------------------")
	for _, b := range sorted {
		fmt.Printf("  %2d | %4d | %s\n", b.ID, b.Year, b.Title)
	}
}
