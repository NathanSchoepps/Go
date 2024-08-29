package models

import (
	"database/sql"
	"example/bookstore/database"
	"fmt"
)

type Book struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

func GetBooks() ([]Book, error) {
	rows, err := database.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func GetBookByID(id string) (Book, error) {
	var book Book
	query := "SELECT * FROM books WHERE id = ?"
	row := database.DB.QueryRow(query, id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
	if err == sql.ErrNoRows {
		return book, fmt.Errorf("Book not found")
	} else if err != nil {
		return book, err
	}
	return book, nil
}

func (b *Book) AddBook() error {
	query := "INSERT INTO books (title, author, price) VALUES (?, ?, ?)"
	_, err := database.DB.Exec(query, b.Title, b.Author, b.Price)
	return err
}
