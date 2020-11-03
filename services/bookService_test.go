package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go-rest-webservices-book-library/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type booksRepositoryMock struct{}

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

func TestIsBookValidFalseForInvalidData(t *testing.T) {
	var jsonStrs = [][]byte{
		[]byte(`{"Name":"Book", "Author": ""}`),
		[]byte(`{"Name", "", "Author": "Author"}`),
		[]byte(`{"Name": "", "Author": ""}`),
		[]byte(`{}`),
	}

	for _, jsonStr := range jsonStrs {
		r, _ := http.NewRequest("", "", bytes.NewBuffer(jsonStr))
		_, valid := isValidData(r)
		_ = assert.False(t, valid)
	}
}

func TestIsBookValidTrueForValidData(t *testing.T) {
	var jsonStr = []byte(`{"Name":"Book", "Author": "Author"}`)
	r, _ := http.NewRequest("", "", bytes.NewBuffer(jsonStr))
	_, valid := isValidData(r)
	_ = assert.True(t, valid)
}

func TestGetString(t *testing.T) {
	var book = domain.Book{
		Id:     1,
		Name:   "Book",
		Author: "Author",
	}

	stringBook := getString(book)
	assert.Equal(t, "{\"Id\":1,\"Name\":\"Book\",\"Author\":\"Author\"}", stringBook)
}

func TestDeleteBooksSuccess(t *testing.T) {
	booksRepository = booksRepositoryMock{}
	booksRepositoryDeleteMock = func(id string) error {
		return nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/book/4", nil)
	vars := map[string]string{
		"id": "4",
	}
	r = mux.SetURLVars(r, vars)

	BookHandler(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteBooksFailure(t *testing.T) {
	booksRepositoryDeleteMock = func(id string) error {
		return errors.New("something bad happened")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/book/4", nil)
	vars := map[string]string{
		"id": "4",
	}
	r = mux.SetURLVars(r, vars)

	BookHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddBookHandlerSuccess(t *testing.T) {
	booksRepositoryAddMock = func(book domain.Book) (int64, error) {
		return 0, nil
	}

	var jsonStr = []byte(`{"Name":"Book","Author":"Author"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")

	AddBookHandler(w, r)

	var book domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), book.Id)
	assert.Equal(t, "Book", book.Name)
	assert.Equal(t, "Author", book.Author)
}

func TestAddBookHandlerFailureBadData(t *testing.T) {
	booksRepositoryAddMock = func(book domain.Book) (int64, error) {
		return 0, nil
	}

	w := httptest.NewRecorder()

	var jsonStrs = [][]byte{
		[]byte(`{"Name":"Book"}`),
		[]byte(`{"Author":"Author"}`),
		[]byte(`{}`),
	}

	for _, jsonStr := range jsonStrs {
		r, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		AddBookHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestAddBookHandlerFailureDatabaseError(t *testing.T) {
	booksRepositoryAddMock = func(book domain.Book) (int64, error) {
		return 0, errors.New("error occurred while inserting record into database")
	}

	var jsonStr = []byte(`{"Name":"Book", "Author":"Author"}`)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")

	AddBookHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBookSuccess(t *testing.T) {
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{{Id: 8, Name: "Book", Author: "Author"}}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)

	BookHandler(w, r)

	var book domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(8), book.Id)
	assert.Equal(t, "Book", book.Name)
	assert.Equal(t, "Author", book.Author)
}

func TestGetBookFailureDatabaseError(t *testing.T) {
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{}, errors.New("error while fetching data")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	BookHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBookNoRecordsFound(t *testing.T) {
	booksRepositoryGetMock = func(id string) ([]domain.Book, error) {
		return []domain.Book{}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/book/1", nil)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	BookHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllBooksHandlerSuccess(t *testing.T) {
	var books = []domain.Book{
		{Id: 1, Name: "Book1", Author: "Author1"},
		{Id: 2, Name: "Book2", Author: "Author2"},
		{Id: 3, Name: "Book3", Author: "Author3"},
	}

	booksRepositoryGetAllMock = func() ([]domain.Book, error) {
		return books, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/books", nil)

	GetAllBooksHandler(w, r)

	var book []domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	assert.Equal(t, http.StatusOK, w.Code)

	for i, book := range books {
		i = i + 1
		index := strconv.FormatInt(int64(i), 10)
		assert.Equal(t, int64(i), book.Id)
		assert.Equal(t, "Book"+index, book.Name)
		assert.Equal(t, "Author"+index, book.Author)
	}
}

func TestGetAllBooksHandlerFailureDatabaseError(t *testing.T) {
	booksRepositoryGetAllMock = func() ([]domain.Book, error) {
		return []domain.Book{}, errors.New("error while getting data from database")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/books", nil)

	GetAllBooksHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateBookSuccess(t *testing.T) {
	booksRepositoryUpdateMock = func(book domain.Book, id string) error {
		return nil
	}

	var jsonStr = []byte(`{"Name":"Book","Author":"Author"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/book/1", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	BookHandler(w, r)

	var book domain.Book
	_ = json.NewDecoder(w.Body).Decode(&book)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), book.Id)
	assert.Equal(t, "Book", book.Name)
	assert.Equal(t, "Author", book.Author)
}

func TestUpdateBookFailureDatabaseError(t *testing.T) {
	booksRepositoryUpdateMock = func(book domain.Book, id string) error {
		return errors.New("error while updating record in database")
	}

	var jsonStr = []byte(`{"Name":"Book","Author":"Author"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/book/1", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	BookHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateBookFailureBadData(t *testing.T) {
	booksRepositoryUpdateMock = func(book domain.Book, id string) error {
		return nil
	}

	w := httptest.NewRecorder()
	var jsonStrs = [][]byte{
		[]byte(`{"Name":"Book"}`),
		[]byte(`{"Author":"Author"}`),
		[]byte(`{}`),
	}

	for _, jsonStr := range jsonStrs {
		r, _ := http.NewRequest("PUT", "/book/1", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		r = mux.SetURLVars(r, map[string]string{"id": "1"})

		BookHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestUnsupportedMethods(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("OPTIONS", "/book/1", nil)

	BookHandler(w, r)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
