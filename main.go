package main

import (
	"context"
	"log"
	"os"
	"time"

	"library/internal/app"
	"library/internal/controller"
	"library/internal/repository"
	"library/internal/story"

	"github.com/gin-gonic/gin"
)

const (
	defaultDSN = "postgres://postgres:postgres@localhost:5434/library?sslmode=disable"
	httpPORT   = ":8081"
)

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
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	bookRepo := repository.NewBookRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	readingRepo := repository.New(pool)
	libraryRepo := repository.NewLibraryRepository(pool)

	lib := story.New(bookRepo, userRepo, readingRepo, libraryRepo)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api, lib)

	if err := r.Run(httpPORT); err != nil {
		log.Fatal(err)
	}
}
