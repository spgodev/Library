package controller

import (
	"errors"
	"library/internal/domain"
	"net/http"
	"strconv"
	"strings"
	"time"

	"library/internal/story"

	"github.com/gin-gonic/gin"
)

type createLibraryRequest struct {
	Name string `json:"name"`
}

func RegisterRoutes(rg *gin.RouterGroup, lib *story.LibraryStory) {
	//1
	rg.POST("/libraries", func(c *gin.Context) {
		var req createLibraryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name_required"})
			return
		}

		created, err := lib.CreateLibrary(c.Request.Context(), req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusCreated, created)
	})
	type createBookRequest struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
		Pages  int    `json:"pages"`
	}
	//2
	rg.POST("/libraries/:libraryId/books", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}

		var req createBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}

		req.Title = strings.TrimSpace(req.Title)
		req.Author = strings.TrimSpace(req.Author)
		if req.Title == "" || req.Author == "" || req.Year <= 0 || req.Pages <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
			return
		}

		created, err := lib.AddBookToLibrary(c.Request.Context(), libraryID, domain.Book{
			Title:  req.Title,
			Author: req.Author,
			Year:   req.Year,
			Pages:  req.Pages,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusCreated, created)
	})
	//3
	type createUserRequest struct {
		FullName string `json:"full_name"`
	}

	rg.POST("/libraries/:libraryId/users", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}

		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}

		req.FullName = strings.TrimSpace(req.FullName)
		if req.FullName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "full_name_required"})
			return
		}

		created, err := lib.AddUserToLibrary(c.Request.Context(), libraryID, req.FullName)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusCreated, created)
	})
	//4,6
	rg.GET("/libraries/:libraryId/books", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}

		var yearPtr *int
		if y := strings.TrimSpace(c.Query("year")); y != "" {
			yi, err := strconv.Atoi(y)
			if err != nil || yi <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bad_year"})
				return
			}
			yearPtr = &yi
		}

		var authorPtr *string
		if a := strings.TrimSpace(c.Query("author")); a != "" {
			authorPtr = &a
		}

		var titlePtr *string
		if t := strings.TrimSpace(c.Query("title")); t != "" {
			titlePtr = &t
		}

		books, err := lib.GetBooksByFilters(c.Request.Context(), libraryID, yearPtr, authorPtr, titlePtr)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusOK, books)
	})
	//5
	rg.GET("/libraries/:libraryId/books/by-author", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}

		author := strings.TrimSpace(c.Query("author"))
		if author == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "author_required"})
			return
		}

		books, err := lib.GetBooksByAuthorInLibrary(c.Request.Context(), libraryID, author)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusOK, books)
	})
	//7
	type addReadingRequest struct {
		BookID   int64  `json:"book_id"`
		ReadDate string `json:"read_date"`
	}

	rg.POST("/libraries/:libraryId/users/:userId/readings", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}
		userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_user_id"})
			return
		}

		var req addReadingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}
		if req.BookID <= 0 || strings.TrimSpace(req.ReadDate) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
			return
		}

		readDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.ReadDate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_date"})
			return
		}

		created, err := lib.AddReading(c.Request.Context(), libraryID, userID, req.BookID, readDate)
		if err != nil {
			if errors.Is(err, domain.NotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
				return
			}
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusCreated, created)
	})
	//8
	rg.GET("/libraries/:libraryId/users/:userId/readings", func(c *gin.Context) {
		libraryID, err := strconv.ParseInt(c.Param("libraryId"), 10, 64)
		if err != nil || libraryID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_library_id"})
			return
		}
		userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_user_id"})
			return
		}

		items, err := lib.GetReadingsByUser(c.Request.Context(), libraryID, userID)
		if err != nil {
			if errors.Is(err, domain.NotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
				return
			}
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		c.JSON(http.StatusOK, items)
	})
}
