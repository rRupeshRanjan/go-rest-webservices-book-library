This is a Golang REST API project, built using  go-v1.15, demonstrating usage of
 
    1. sql database
    2. rest apis
    3. mux router (for routing and limiting different request methods)
    4. unit and benchmarking tests

This project has below functionalities:

    1. See all books available in library
    2. Add a book to library
    3. Update a book in library (using its id)
    4. Delete a book from library (using id)

Data is being stored in SQL table named "books" to store books, with columns named as below, which are the book attributes

    1. unique id
    2. Name
    3. Author 

#### How to run test cases

1. "testing" library for creating tests
2. Test files nomenclature:
    - filename must have _test as suffix (eg. bookservice_test.go)
    - tests **MUST** start with Test (e.g. func TestSomeHandler())
    - tests **MUST** take input as pointer of testing.T (e.g. TestSomeHandler(t *testing.T)) 
    - benchmark tests **MUST** start with Benchmark word (e.g. func BenchmarkSomeHandler())
    - benchmark tests **MUST** take input as pointer of testing.B (e.g. TestSomeHandler(t *testing.B)) 
3. To run tests for specific package
    - Switch to respective directory
    - go test
4. More for tests
    - go test -v (for verbose)
    - go test ./... ( run all tests) 
    - go test -bench . (for running benchmark tests)
    - go test -bench . -benchtime 10s (specify time for each benchmark test)
    - go test -bench . -benchmem (for memory profiling of benchmark tests)
    - go test -coverprofile cover.txt (Run test with coverage)
    - go tool cover -html cover.txt (to see to html for coverage file)
    - go test -coverprofile count.out -covermode count (gives relative coverage for code)
