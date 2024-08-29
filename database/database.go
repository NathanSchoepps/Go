package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
	dsn := "root:@tcp(localhost:3306)/book"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	fmt.Println("Connected to the database succesfully")
	return nil
}
