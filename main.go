package main

import (
	"book-library/repository"
	"book-library/services"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	repository.InitBooksDb()
	router := mux.NewRouter()

	router.HandleFunc("/books", services.GetAllBooksHandler).Methods("GET")
	router.HandleFunc("/book", services.AddBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", services.BookHandler).Methods("GET", "DELETE", "PUT")

	_ = http.ListenAndServe(":8080", router)
}
