package services

import (
	"bytes"
	"github.com/gorilla/mux"
	"go-rest-webservices-book-library/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkSetup(b *testing.B) {
	booksRepository = booksRepositoryMock{}
}

func BenchmarkIsValidBook(b *testing.B) {
	data := []byte(`{"Name":"Book", "Author": "Author"}`)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r, _ := http.NewRequest("", "", bytes.NewBuffer(data))
		isValidBook(r)
	}
}

func BenchmarkGetString(b *testing.B) {
	data := domain.Book{Name: "Book", Author: "Author"}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getString(data)
	}
}

func BenchmarkBookHandlerGetBookById(b *testing.B) {
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{{Id: 8, Name: "Book", Author: "Author"}}, nil
	}
	r, _ := http.NewRequest("GET", "/book/8", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "8"})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		BookHandler(w, r)
	}
}

func BenchmarkBookHandlerUpdateBook(b *testing.B) {
	booksRepositoryUpdateMock = func(book domain.Book, id string) error {
		return nil
	}
	data := []byte(`{"Name":"Book", "Author":"Author"}`)
	r, _ := http.NewRequest("PUT", "/book/1", bytes.NewBuffer(data))
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r, _ = http.NewRequest("PUT", "/book/1", bytes.NewBuffer(data))
		BookHandler(w, r)
	}
}

func BenchmarkBookHandlerDeleteBook(b *testing.B) {
	booksRepositoryDeleteMock = func(id string) error {
		return nil
	}
	r, _ := http.NewRequest("DELETE", "/book/4", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "4"})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		BookHandler(w, r)
	}
}

func BenchmarkBookHandlerUnsupportedMethods(b *testing.B) {
	r, _ := http.NewRequest("OPTIONS", "/book/1", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		BookHandler(w, r)
	}
}

func BenchmarkGetAllBooksHandler(b *testing.B) {
	booksRepositoryGetAllMock = func() ([]domain.Book, error) {
		return []domain.Book{
			{Id: 1, Name: "Book1", Author: "Author1"},
			{Id: 2, Name: "Book2", Author: "Author2"},
			{Id: 3, Name: "Book3", Author: "Author3"},
		}, nil
	}
	r, _ := http.NewRequest("GET", "/books", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		GetAllBooksHandler(w, r)
	}
}

func BenchmarkAddBookHandler(b *testing.B) {
	booksRepositoryAddMock = func(book domain.Book) (int64, error) {
		return 0, nil
	}
	data := []byte(`{"Name":"Book", "Author":"Author"}`)
	r, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(data))
	r.Header.Set("Content-Type", "application/json")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r, _ = http.NewRequest("POST", "/book", bytes.NewBuffer(data))
		AddBookHandler(w, r)
	}
}
