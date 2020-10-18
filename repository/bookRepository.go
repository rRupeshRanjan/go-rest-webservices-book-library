package repository

import (
	"book-library/domain"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var database *sql.DB

const (
	createQuery = "SELECT * FROM books WHERE id=?"
	updateQuery = "UPDATE books SET name=?, author=? where id=?"
	deleteQuery = "DELETE FROM books WHERE id=?"
	getAllQuery = "SELECT id, name, author FROM books"
	insertQuery = "INSERT INTO books (name, author) VALUES (?, ?)"
	initializeDatabaseQuery =
		`CREATE TABLE IF NOT EXISTS books (
					id INTEGER PRIMARY KEY, 
					name TEXT, 
					author TEXT);`
)

func InitBooksDb() {
	database, _ = sql.Open("sqlite3", "books.sql")
	_, err := database.Exec(initializeDatabaseQuery)
	if err != nil {
		log.Panic("Failure while initializing database, {}", err.Error())
	}
}

func UpdateBook(book domain.Book, id string) error {
	statement, _ := database.Prepare(updateQuery)
	_, err := statement.Exec(book.Name, book.Author, id)

	return err
}

func GetBook(id string) (*sql.Rows, error) {
	rows, err := database.Query(createQuery, id)
	return rows, err
}

func DeleteBook(id string) error {
	statement, _ := database.Prepare(deleteQuery)
	_, err := statement.Exec(id)

	return err
}

func GetAllBooks() *sql.Rows {
	rows, _ := database.Query(getAllQuery)
	return rows
}

func AddBook(book domain.Book) (int64, error) {
	statement, _ := database.Prepare(insertQuery)
	result, insertRecordErr := statement.Exec(book.Name, book.Author)
	if insertRecordErr != nil {
		log.Printf("Error occured while inserting data in books table: %s", insertRecordErr)
	}
	return result.LastInsertId()
}
