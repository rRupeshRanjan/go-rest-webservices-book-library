package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-rest-webservices-book-library/config"
	"go-rest-webservices-book-library/domain"
	"go.uber.org/zap"
)

var (
	database *sql.DB
	logger   *zap.Logger
)

const (
	getQuery                = "SELECT * FROM books WHERE id=?"
	updateQuery             = "UPDATE books SET name=?, author=? where id=?"
	deleteQuery             = "DELETE FROM books WHERE id=?"
	getAllQuery             = "SELECT * FROM books"
	insertQuery             = "INSERT INTO books (name, author) VALUES (?, ?)"
	initializeDatabaseQuery = `CREATE TABLE IF NOT EXISTS books (
									id INTEGER PRIMARY KEY, 
									name TEXT, 
									author TEXT);`
)

func init() {
	initBooksDb()
	logger = config.AppLogger
}

func initBooksDb() {
	database, _ = sql.Open("sqlite3", "books.sql")
	_, err := database.Exec(initializeDatabaseQuery)
	if err != nil {
		logger.Panic("Failure while initializing database, {}" + err.Error())
	}
}

func UpdateBook(book domain.Book, id string) error {
	statement, _ := database.Prepare(updateQuery)
	_, err := statement.Exec(book.Name, book.Author, id)

	return err
}

func GetBook(id string) ([]domain.Book, error) {
	rows, err := database.Query(getQuery, id)
	var books []domain.Book

	if err == nil && rows.Next() {
		var book domain.Book
		_ = rows.Scan(&book.Id, &book.Name, &book.Author)
		books = append(books, book)
	}

	return books, err
}

func DeleteBook(id string) error {
	statement, _ := database.Prepare(deleteQuery)
	_, err := statement.Exec(id)

	return err
}

func GetAllBooks() ([]domain.Book, error) {
	rows, err := database.Query(getAllQuery)

	var books []domain.Book

	for err == nil && rows.Next() {
		var book domain.Book
		err = rows.Scan(&book.Id, &book.Name, &book.Author)
		books = append(books, book)
	}

	return books, err
}

func AddBook(book domain.Book) (int64, error) {
	statement, _ := database.Prepare(insertQuery)
	result, insertRecordErr := statement.Exec(book.Name, book.Author)
	if insertRecordErr != nil {
		logger.Error("Error occurred while inserting data in books table: %s" + insertRecordErr.Error())
		return -1, insertRecordErr
	}
	return result.LastInsertId()
}
