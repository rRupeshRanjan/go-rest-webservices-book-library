package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go-rest-webservices-book-library/domain"
	"go-rest-webservices-book-library/repository"
	"log"
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

var booksRepository BooksRepositoryInterface

func Init() {
	repository.InitBooksDb()
	booksRepository = BooksRepository{}
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
	}
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book domain.Book
	requestBody := r.Body
	decodeErr := json.NewDecoder(requestBody).Decode(&book)

	if decodeErr == nil && isValidData(book) {
		updateErr := booksRepository.updateBook(book, id)
		if updateErr == nil {
			book.Id, _ = strconv.ParseInt(id, 10, 64)
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, getString(book))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		log.Printf("Improper data passed for update: %s", getString(book))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	books, err := booksRepository.getBook(id)

	if err == nil {
		if len(books) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, getString(books[0]))
		}
	} else {
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
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetAllBooksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := booksRepository.getAllBooks()

	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, getString(books))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book domain.Book
	requestBody := r.Body
	decodeErr := json.NewDecoder(requestBody).Decode(&book)

	if decodeErr == nil && isValidData(book) {
		rowId, insertRecordErr := booksRepository.addBook(book)
		if insertRecordErr == nil {
			book.Id = rowId
			_, _ = fmt.Fprintf(w, getString(book))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		log.Printf("Improper data passed for update: %s", getString(book))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func isValidData(book domain.Book) bool {
	return book.Name != "" && book.Author != ""
}

func getString(input interface{}) string {
	jsonDeserializedObject, deserializationErr := json.Marshal(input)
	stringObject := ""

	if deserializationErr == nil {
		stringObject = string(jsonDeserializedObject)
	} else {
		log.Printf("Error while deserializing data: %s", deserializationErr)
	}

	return stringObject
}
