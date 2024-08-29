package handlers

import (
	"example/bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books, err := models.GetBooks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	book, err := models.GetBookByID(id)
	if err != nil {
		if err.Error() == "book not found" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func PostBooks(c *gin.Context) {
	var newBook models.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	if err := newBook.AddBook(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newBook)
}
