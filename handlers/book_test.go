package handlers

import (
	"bytes"
	"encoding/json"
	"example/bookstore/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock de la couche modèle pour isoler les tests
type MockModel struct {
	mock.Mock
}

func (m *MockModel) GetBooks() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockModel) GetBookByID(id string) (models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockModel) AddBook(b *models.Book) error {
	args := m.Called(b)
	return args.Error(0)
}

func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockModel := new(MockModel)
	mockBooks := []models.Book{
		{ID: 1, Title: "Book One", Author: "Author One", Price: 10.99},
		{ID: 2, Title: "Book Two", Author: "Author Two", Price: 12.99},
	}

	mockModel.On("GetBooks").Return(mockBooks, nil)

	router := gin.Default()
	c := router.Group("/")
	c.GET("/books", func(c *gin.Context) {
		books, err := mockModel.GetBooks()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, books)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var books []models.Book
	err := json.Unmarshal(w.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.Equal(t, mockBooks, books)
}

func TestGetBookByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockModel := new(MockModel)
	mockBook := models.Book{ID: 1, Title: "Book One", Author: "Author One", Price: 10.99}

	mockModel.On("GetBookByID", "1").Return(mockBook, nil)

	router := gin.Default()
	c := router.Group("/")
	c.GET("/books/:id", func(c *gin.Context) {
		id := c.Param("id")
		book, err := mockModel.GetBookByID(id)
		if err != nil {
			if err.Error() == "book not found" {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			return
		}
		c.IndentedJSON(http.StatusOK, book)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var book models.Book
	err := json.Unmarshal(w.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, mockBook, book)
}

func TestPostBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockModel := new(MockModel)
	newBook := models.Book{Title: "New Book", Author: "New Author", Price: 15.99}

	mockModel.On("AddBook", &newBook).Return(nil)

	router := gin.Default()
	c := router.Group("/")
	c.POST("/books", func(c *gin.Context) {
		var newBook models.Book
		if err := c.BindJSON(&newBook); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
			return
		}

		if err := mockModel.AddBook(&newBook); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newBook)
	})

	body, _ := json.Marshal(newBook)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBook models.Book
	err := json.Unmarshal(w.Body.Bytes(), &createdBook)
	assert.NoError(t, err)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Price, createdBook.Price)
}
