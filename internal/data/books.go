package data

import (
	"errors"

	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// In-memory store for books
var books = []models.Book{
	{ID: 1, Title: "Blue Train", Author: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Author: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Author: "Sarah Vaughan", Price: 39.99},
}

// GetAllBooks returns all books
func GetAllBooks() []models.Book {
	return books
}

// GetBookByID returns a book by ID
func GetBookByID(id int) (*models.Book, error) {
	for _, b := range books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("book not found")
}

// AddBook adds a new book
func AddBook(b models.Book) models.Book {
	// Generate a new ID
	var maxID int // it gets the zero value of its type by default.
	for _, book := range books {
		if book.ID > maxID {
			maxID = book.ID
		}
	}
	b.ID = maxID + 1
	books = append(books, b)
	return b
}

// UpdateBook updates a book by ID
func UpdateBook(id int, updated models.Book) (*models.Book, error) {
	for i, b := range books {
		if b.ID == id {
			if updated.Title != "" {
				books[i].Title = updated.Title
			}
			if updated.Author != "" {
				books[i].Author = updated.Author
			}
			if updated.Price != 0 {
				books[i].Price = updated.Price
			}
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// DeleteBook deletes a book by ID
func DeleteBook(id int) error {
	for i, b := range books {
		if b.ID == id { // ... means: expand this slice into individual elements
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}
