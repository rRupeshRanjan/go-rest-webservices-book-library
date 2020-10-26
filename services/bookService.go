package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-rest-webservices-book-library/domain"
	"go-rest-webservices-book-library/repository"
	"log"
	"net/http"
)

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
	_ = json.NewDecoder(requestBody).Decode(&book)

	updateErr := repository.UpdateBook(book, id)
	if updateErr == nil {
		_, _ = fmt.Fprintf(w, getString(book))
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Improper data passed for update: %s", requestBody)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rows, err := repository.GetBook(id)

	if err == nil && rows.Next() {
		var book domain.Book
		var id int64
		var name string
		var author string
		_ = rows.Scan(&id, &name, &author)
		book = domain.Book{Id: id, Name: name, Author: author}
		_, _ = fmt.Fprintf(w, getString(book))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deleteError := repository.DeleteBook(id)

	if deleteError == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
	}
}

func GetAllBooksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []domain.Book
	rows := repository.GetAllBooks()

	var id int64
	var name string
	var author string
	for rows.Next() {
		_ = rows.Scan(&id, &name, &author)
		books = append(books, domain.Book{Id: id, Name: name, Author: author})
	}
	_, _ = fmt.Fprintf(w, getString(books))
}

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book domain.Book
	requestBody := r.Body
	decodeErr := json.NewDecoder(requestBody).Decode(&book)

	if decodeErr == nil && isValidData(book) {
		rowId, insertRecordErr := repository.AddBook(book)
		if insertRecordErr == nil {
			book.Id = rowId
			_, _ = fmt.Fprintf(w, getString(book))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		log.Printf("Improper data passed for update: %s", requestBody)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func isValidData(book domain.Book) bool {
	return book.Name != "" && book.Author != ""
}

func getString(input interface{}) string {
	jsonDeserializedObject, deserializationErr := json.Marshal(input)

	if deserializationErr == nil {
		return string(jsonDeserializedObject)
	}
	log.Printf("Error while deserializing data: %s", deserializationErr)
	return ""
}
