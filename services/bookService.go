package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go-rest-webservices-book-library/config"
	"go-rest-webservices-book-library/domain"
	"go-rest-webservices-book-library/repository"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type BooksRepository struct{}

type BooksRepositoryInterface interface {
	getBook(id string) ([]domain.Book, error)
	getAllBooks() ([]domain.Book, error)
	addBook(book domain.Book) (int64, error)
	updateBook(book domain.Book, id string) error
	deleteBook(id string) error
}

var (
	booksRepository BooksRepositoryInterface
	logger          *zap.Logger
)

func init() {
	booksRepository = BooksRepository{}
	logger = config.AppLogger
}

func (b BooksRepository) getBook(id string) ([]domain.Book, error) {
	return repository.GetBook(id)
}

func (b BooksRepository) getAllBooks() ([]domain.Book, error) {
	return repository.GetAllBooks()
}

func (b BooksRepository) addBook(book domain.Book) (int64, error) {
	return repository.AddBook(book)
}

func (b BooksRepository) updateBook(book domain.Book, id string) error {
	return repository.UpdateBook(book, id)
}

func (b BooksRepository) deleteBook(id string) error {
	return repository.DeleteBook(id)
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		getBookHandler(w, r)
	case "DELETE":
		deleteBookHandler(w, r)
	case "PUT":
		updateBookHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	book, valid := isValidBook(r)

	if valid {
		updateErr := booksRepository.updateBook(book, id)
		if updateErr == nil {
			book.Id, _ = strconv.ParseInt(id, 10, 64)
			logger.Info("Successfully updated book: " + getString(book))
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, getString(book))
		} else {
			logger.Error("Error while updating book: " + id + " with error: " + updateErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		logger.Error("Improper data passed for update: " + getString(book))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	books, getBookErr := booksRepository.getBook(id)

	if getBookErr == nil {
		if len(books) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, getString(books[0]))
		}
	} else {
		logger.Error("Error while getting book: " + id + " with error: " + getBookErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deleteError := booksRepository.deleteBook(id)

	if deleteError == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		logger.Error("Error while deleting book: " + id + " with error: " + deleteError.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, getAllError := booksRepository.getAllBooks()

	if getAllError == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, getString(books))
	} else {
		logger.Error("Error while getting all books with error: " + getAllError.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	book, valid := isValidBook(r)

	if valid {
		rowId, insertRecordErr := booksRepository.addBook(book)
		if insertRecordErr == nil {
			book.Id = rowId
			_, _ = fmt.Fprintf(w, getString(book))
		} else {
			logger.Error("Error while creating book with error: " + insertRecordErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		logger.Error("Improper data passed for create: " + getString(book))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func isValidBook(r *http.Request) (domain.Book, bool) {
	var book domain.Book
	decodeErr := json.NewDecoder(r.Body).Decode(&book)

	return book, decodeErr == nil && book.Name != "" && book.Author != ""
}

func getString(input interface{}) string {
	jsonDeserializedObject, deserializationErr := json.Marshal(input)

	if deserializationErr == nil {
		return string(jsonDeserializedObject)
	} else {
		logger.Error("Error while deserializing data:" + deserializationErr.Error())
		return ""
	}
}
