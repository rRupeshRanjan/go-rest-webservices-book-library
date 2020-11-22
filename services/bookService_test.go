package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go-rest-webservices-book-library/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type booksRepositoryMock struct{}

type scenario struct {
	name           string
	book           domain.Book
	books          []domain.Book
	data           []byte
	status         int
	err            error
	method         string
	valid          bool
	expectedString string
}

var (
	booksRepositoryGetMock    func(id string) ([]domain.Book, error)
	booksRepositoryGetAllMock func() ([]domain.Book, error)
	booksRepositoryAddMock    func(book domain.Book) (int64, error)
	booksRepositoryUpdateMock func(book domain.Book, id string) error
	booksRepositoryDeleteMock func(id string) error
)

func (b booksRepositoryMock) getBook(id string) ([]domain.Book, error) {
	return booksRepositoryGetMock(id)
}

func (b booksRepositoryMock) getAllBooks() ([]domain.Book, error) {
	return booksRepositoryGetAllMock()
}

func (b booksRepositoryMock) addBook(book domain.Book) (int64, error) {
	return booksRepositoryAddMock(book)
}

func (b booksRepositoryMock) updateBook(book domain.Book, id string) error {
	return booksRepositoryUpdateMock(book, id)
}

func (b booksRepositoryMock) deleteBook(id string) error {
	return booksRepositoryDeleteMock(id)
}

func TestSetup(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	initLogger()
}

func TestIsValidData(t *testing.T) {
	t.Parallel()
	scenarios := []scenario{
		{
			name:  "Invalid if author name is empty string",
			data:  []byte(`{"Name":"Book", "Author": ""}`),
			valid: false,
		},
		{
			name:  "Invalid if book name is empty string",
			data:  []byte(`{"Name", "", "Author": "Author"}`),
			valid: false,
		},
		{
			name:  "Invalid if book name or author name is empty string",
			data:  []byte(`{"Name": "", "Author": ""}`),
			valid: false,
		},
		{
			name:  "Invalid if fields are missing",
			data:  []byte(`{}`),
			valid: false,
		},
		{
			name:  "Valid when book name and author name and non-empty strings",
			data:  []byte(`{"Name":"Book", "Author": "Author"}`),
			valid: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			r, _ := http.NewRequest("", "", bytes.NewBuffer(scenario.data))
			_, valid := isValidBook(r)
			if scenario.valid != valid {
				t.Errorf("Expected %v, found %v\n", scenario.valid, valid)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()
	scenarios := []scenario{
		{
			name:           "Convert book to string",
			book:           domain.Book{Id: 1, Name: "Book", Author: "Author"},
			expectedString: "{\"Id\":1,\"Name\":\"Book\",\"Author\":\"Author\"}",
		},
		{
			name:           "Convert book to string",
			book:           domain.Book{Id: 1, Name: "Book"},
			expectedString: "{\"Id\":1,\"Name\":\"Book\",\"Author\":\"\"}",
		},
		{
			name:           "Convert book to string",
			book:           domain.Book{Name: "Book", Author: "Author"},
			expectedString: "{\"Id\":0,\"Name\":\"Book\",\"Author\":\"Author\"}",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			got := getString(scenario.book)
			if scenario.expectedString != got {
				t.Errorf("Expected %v, got %v\n", scenario.expectedString, got)
			}
		})
	}
}

func TestDeleteBooks(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("DELETE", "/book/4", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "4"})

	scenarios := []scenario{
		{
			name:   "should delete with status code 204",
			status: http.StatusNoContent,
		},
		{
			name:   "should give 500 on delete error",
			err:    errors.New("something bad happened"),
			status: http.StatusInternalServerError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			booksRepositoryDeleteMock = func(id string) error {
				return scenario.err
			}
			BookHandler(w, r)
			if w.Code != scenario.status {
				t.Errorf("Expected status code: %v, got %v", scenario.status, w.Code)
			}
		})
	}
}

func TestAddBookHandler(t *testing.T) {
	t.Parallel()
	scenarios := []scenario{
		{
			name:   "success for book create",
			book:   domain.Book{Name: "Book", Author: "Author"},
			data:   []byte(`{"Name":"Book","Author":"Author"}`),
			status: http.StatusOK,
		},
		{
			name:   "failure for book create for bad data",
			data:   []byte(`{"Name":"Book"`),
			status: http.StatusBadRequest,
		},
		{
			name:   "failure for book create for bad data",
			data:   []byte(`{"Author":"Author"}`),
			status: http.StatusBadRequest,
		},
		{
			name:   "failure for book create for bad data",
			data:   []byte(`{}`),
			status: http.StatusBadRequest,
		},
		{
			name:   "failure for book create for db errors",
			err:    errors.New("error occurred while inserting record into database"),
			data:   []byte(`{"Name":"Book", "Author":"Author"}`),
			status: http.StatusInternalServerError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			booksRepositoryAddMock = func(book domain.Book) (int64, error) {
				return 0, scenario.err
			}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(scenario.data))
			r.Header.Set("Content-Type", "application/json")
			AddBookHandler(w, r)

			compareResponses(t, w, scenario)
		})
	}
}

func TestGetBookByIdHandler(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/book/8", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "8"})
	scenarios := []scenario{
		{
			name:   "should successfully get book by id",
			books:  []domain.Book{{Id: 8, Name: "Book", Author: "Author"}},
			status: http.StatusOK,
		},
		{
			name:   "should give 404 for get book by id",
			books:  []domain.Book{},
			status: http.StatusNotFound,
		},
		{
			name:   "should give 500 for get book by id for database errors",
			books:  []domain.Book{},
			err:    errors.New("error while fetching data"),
			status: http.StatusInternalServerError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
				return scenario.books, scenario.err
			}
			BookHandler(w, r)

			if w.Code != scenario.status {
				t.Errorf("Expected status code: %v, Got: %v", scenario.status, w.Code)
			}

			if w.Code == http.StatusOK {
				var book domain.Book
				_ = json.NewDecoder(w.Body).Decode(&book)

				if book != scenario.books[0] {
					t.Errorf("Expected Data: %v, Got: %v", scenario.books[0], book)
				}
			}
		})
	}
}

func TestGetAllBooksHandler(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/books", nil)
	scenarios := []scenario{
		{
			name: "should get all books",
			books: []domain.Book{
				{Id: 1, Name: "Book1", Author: "Author1"},
				{Id: 2, Name: "Book2", Author: "Author2"},
				{Id: 3, Name: "Book3", Author: "Author3"},
			},
			status: http.StatusOK,
		},
		{
			name:   "should give 500 for getAll books for database error",
			books:  []domain.Book{},
			status: http.StatusInternalServerError,
			err:    errors.New("error while getting data from database"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			booksRepositoryGetAllMock = func() ([]domain.Book, error) {
				return scenario.books, scenario.err
			}
			GetAllBooksHandler(w, r)
			compareResponses(t, w, scenario)
		})
	}
}

func TestUpdateBookHandler(t *testing.T) {
	t.Parallel()
	scenarios := []scenario{
		{
			name:   "should update record",
			book:   domain.Book{Id: 1, Name: "Book", Author: "Author"},
			data:   []byte(`{"Name":"Book","Author":"Author"}`),
			status: http.StatusOK,
		},
		{
			name:   "should fail update record for database errors",
			data:   []byte(`{"Name":"Book","Author":"Author"}`),
			err:    errors.New("error while updating record in database"),
			status: http.StatusInternalServerError,
		},
		{
			name:   "should fail update record for bad data",
			data:   []byte(`{"Name":"Book"}`),
			status: http.StatusBadRequest,
		},
		{
			name:   "should fail update record for bad data",
			data:   []byte(`{"Author":"Author"}`),
			status: http.StatusBadRequest,
		},
		{
			name:   "should fail update record for bad data",
			data:   []byte(`{}`),
			status: http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			booksRepositoryUpdateMock = func(book domain.Book, id string) error {
				return scenario.err
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("PUT", "/book/1", bytes.NewBuffer(scenario.data))
			r.Header.Set("Content-Type", "application/json")
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			BookHandler(w, r)
			compareResponses(t, w, scenario)
		})
	}
}

func TestUnsupportedMethods(t *testing.T) {
	t.Parallel()
	scenarios := []scenario{
		{
			name:   "options method not supported",
			method: "OPTIONS",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "head method not supported",
			method: "HEAD",
			status: http.StatusMethodNotAllowed,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(scenario.method, "/book/1", nil)
			BookHandler(w, r)

			if scenario.status != w.Code {
				t.Errorf("EXpected status code: %v, got %v", scenario.status, w.Code)
			}
		})
	}
}

func compareResponses(t *testing.T, w *httptest.ResponseRecorder, scenario scenario) {
	if w.Code != scenario.status {
		t.Errorf("Expected status code: %v, Got: %v", scenario.status, w.Code)
	}

	if w.Code == http.StatusOK {
		if len(scenario.books) == 0 {
			var book domain.Book
			_ = json.NewDecoder(w.Body).Decode(&book)

			var books []domain.Book
			_ = json.NewDecoder(w.Body).Decode(&books)

			if book != scenario.book || len(books) != len(scenario.books) {
				t.Errorf("Expected %v, got %v", scenario.book, book)
			}
		} else {
			var books []domain.Book
			_ = json.NewDecoder(w.Body).Decode(&books)

			if len(books) != len(scenario.books) {
				t.Errorf("Expected %v, got %v", scenario.books, books)
			}
		}
	}
}
