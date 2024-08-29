package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example/bookstore/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetBooks tests the GetBooks handler
func TestGetBooks(t *testing.T) {
	// Setup Gin router
	router := gin.Default()
	router.GET("/books", GetBooks)

	// Mock the GetBooks function in models
	models.GetBooks = func() ([]models.Book, error) {
		return []models.Book{
			{ID: 1, Title: "1984", Author: "George Orwell", Price: 9.99},
			{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 7.99},
		}, nil
	}

	// Create a request to send to the endpoint
	req, _ := http.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1984")
	assert.Contains(t, w.Body.String(), "To Kill a Mockingbird")
}

// TestGetBookByID tests the GetBookByID handler
func TestGetBookByID(t *testing.T) {
	// Setup Gin router
	router := gin.Default()
	router.GET("/books/:id", GetBookByID)

	// Mock the GetBookByID function in models
	models.GetBookByID = func(id string) (models.Book, error) {
		if id == "1" {
			return models.Book{ID: 1, Title: "1984", Author: "George Orwell", Price: 9.99}, nil
		}
		return models.Book{}, fmt.Errorf("book not found")
	}

	// Test for an existing book
	req, _ := http.NewRequest(http.MethodGet, "/books/1", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1984")

	// Test for a non-existent book
	req, _ = http.NewRequest(http.MethodGet, "/books/999", nil)
	w = httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Book not found")
}

// TestPostBooks tests the PostBooks handler
func TestPostBooks(t *testing.T) {
	// Setup Gin router
	router := gin.Default()
	router.POST("/books", PostBooks)

	// Mock the AddBook function in models
	models.AddBook = func(book *models.Book) error {
		book.ID = 3 // Assume ID is auto-incremented to 3
		return nil
	}

	// Create a request to send to the endpoint
	body := `{"title": "Brave New World", "author": "Aldous Huxley", "price": 8.99}`
	req, _ := http.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Brave New World")
}
