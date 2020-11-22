package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go-rest-webservices-book-library/config"
	"go-rest-webservices-book-library/services"
	"go-rest-webservices-book-library/startup"
	"net/http"
)

func main() {
	startup.Initialize()

	logFile := config.LogFile
	router := mux.NewRouter()

	router.Handle(
		"/books",
		handlers.LoggingHandler(logFile, http.HandlerFunc(services.GetAllBooksHandler))).
		Methods("GET")

	router.Handle(
		"/book",
		handlers.LoggingHandler(logFile, http.HandlerFunc(services.AddBookHandler))).
		Methods("POST")

	router.Handle(
		"/book/{id}",
		handlers.LoggingHandler(logFile, http.HandlerFunc(services.BookHandler))).
		Methods("GET", "DELETE", "PUT")

	_ = http.ListenAndServe(config.ServerPort, router)
}
