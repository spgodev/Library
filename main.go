package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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

	if err := app.RunMigrations(dsn); err != nil {
		log.Fatal(err)
	}

	pool, err := app.Open(ctx, dsn)
	defer pool.Close()

	bookRepo := repository.NewBookRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	readingRepo := repository.New(pool)

	lib := story.New(bookRepo, userRepo, readingRepo)

	ok, err := lib.HasBookTitle(ctx, "1984")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Has '1984':", ok)

	found, err := lib.FindBookByTitle(ctx, "Design Patterns")
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Find by title = Design Patterns", found)

	readers, err := lib.GetReadersByBook(ctx, int64(1))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nReaders of book #1:")
	for _, r := range readers {
		fmt.Printf("  - %-12s %s\n", r.User, r.Date.Format("2006-01-02"))
	}

	br, err := lib.MarkAsRead(ctx, 3, "Vasya Pupkin", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	fmt.Println("\nMarkAsRead:", br, "err:", err)

	byYear, err := lib.GetBooksByYear(ctx, 1949)
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Books by year = 1949", byYear)

	byAuthor, err := lib.GetBooksByAuthor(ctx, "George Orwell")
	if err != nil {
		log.Fatal(err)
	}
	printBooks("Books by author = George Orwell", byAuthor)

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
