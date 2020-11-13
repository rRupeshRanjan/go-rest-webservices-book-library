package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go-rest-webservices-book-library/services"
	"net/http"
	"os"
)

func main() {
	file, _ := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	services.Init()
	router := mux.NewRouter()

	router.Handle(
		"/books",
		handlers.LoggingHandler(file, http.HandlerFunc(services.GetAllBooksHandler))).
		Methods("GET")

	router.Handle(
		"/book",
		handlers.LoggingHandler(file, http.HandlerFunc(services.AddBookHandler))).
		Methods("POST")

	router.Handle(
		"/book/{id}",
		handlers.LoggingHandler(file, http.HandlerFunc(services.BookHandler))).
		Methods("GET", "DELETE", "PUT")

	_ = http.ListenAndServe(":8080", router)
}
