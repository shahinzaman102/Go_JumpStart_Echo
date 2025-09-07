package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// GetBooks returns all books in JSON format.
func GetBooks(c echo.Context) error {
	books := data.GetAllBooks()
	return c.JSON(http.StatusOK, books)
}

// GetBookByID returns a single book by ID.
func GetBookByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid book ID"})
	}

	book, err := data.GetBookByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "book not found"})
	}

	return c.JSON(http.StatusOK, book)
}

// PostBook creates a new book.
func PostBook(c echo.Context) error {
	var newBook models.Book
	if err := c.Bind(&newBook); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid JSON"})
	}

	// Trim whitespace
	newBook.Title = strings.TrimSpace(newBook.Title)
	newBook.Author = strings.TrimSpace(newBook.Author)

	// Validation
	if newBook.Title == "" || newBook.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title and Author are required"})
	}
	if len(newBook.Title) > 200 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title must be 1-200 characters"})
	}
	if len(newBook.Author) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Author must be 1-100 characters"})
	}

	book := data.AddBook(newBook)
	return c.JSON(http.StatusCreated, book)
}

// UpdateBook updates an existing book by ID.
func UpdateBook(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid book ID"})
	}

	var updatedData models.Book
	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid JSON"})
	}

	// Trim whitespace
	updatedData.Title = strings.TrimSpace(updatedData.Title)
	updatedData.Author = strings.TrimSpace(updatedData.Author)

	// Validate at least one field
	if updatedData.Title == "" && updatedData.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No fields to update"})
	}

	if updatedData.Title != "" && len(updatedData.Title) > 200 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title must be 1-200 characters"})
	}
	if updatedData.Author != "" && len(updatedData.Author) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Author must be 1-100 characters"})
	}

	updatedBook, err := data.UpdateBook(id, updatedData)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "book not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  "success",
		"message": "book updated successfully",
		"book":    updatedBook,
	})
}

// DeleteBook removes a book by ID.
func DeleteBook(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid book ID"})
	}

	if err := data.DeleteBook(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "book not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "book deleted successfully",
	})
}
