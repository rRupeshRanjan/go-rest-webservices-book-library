This is a Golang REST API project, built using  go-v1.15, sql database, and mux.

This has below functionalities:
1. See all books available in library
2. Add a book to library
3. Update a book in library (using its id)
4. Delete a book from library (using id)

Every book has below attributes:
1. unique id
2. Name
3. Author 

Data stored in SQL table named "books" to store books, with columns named as above.

Mux is used for request routing, and limiting the methods allowed for a particular API.

How to run Test cases in Golang
1. "testing" library for creating tests
2. To run tests
    - Switch to respective directory
    - go test
    - go test -v (for verbose)
    - go test -coverprofile=cover.txt (Run test with coverage)
    - go tool cover -html=cover.txt -o cover.html (to convert txt to html for coverage file)
