package module

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tira.com/src/db"
	"tira.com/src/helper"
	"tira.com/src/model"

	"github.com/gorilla/mux"
)

// Get all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")

	db := db.SetupDBPostgres()

	helper.PrintMessage("Getting books...")

	// Get all books from books table that don't have bookID = "1"
	rows, err := db.Query("SELECT * FROM books where bookID <> $1", "1")

	helper.CheckErr(err)
	var books []model.Book
	// var response []JsonResponse
	// Foreach book
	for rows.Next() {
		var id int
		var bookID string
		var bookName string

		err = rows.Scan(&id, &bookID, &bookName)

		helper.CheckErr(err)

		books = append(books, model.Book{BookID: bookID, BookName: bookName})
	}

	var response = model.JsonResponse{Type: "success", Data: books}

	json.NewEncoder(w).Encode(response)
}

// Create a book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "POST")

	// Declare a new Book struct.
	var p model.Book

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookID := p.BookID
	bookName := p.BookName

	var response = model.JsonResponse{}

	if bookID == "" || bookName == "" {
		response = model.JsonResponse{Type: "error", Message: "You are missing bookID or bookName parameter."}
	} else {
		db := db.SetupDBPostgres()

		helper.PrintMessage("Inserting book into DB")

		fmt.Println("Inserting new book with ID: " + bookID + " and name: " + bookName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO books(bookID, bookName) VALUES($1, $2) returning id;", bookID, bookName).Scan(&lastInsertID)

		helper.CheckErr(err)

		response = model.JsonResponse{Type: "success", Message: "The book has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	params := mux.Vars(r)

	bookID := params["bookid"]

	var response = model.JsonResponse{}

	if bookID == "" {
		response = model.JsonResponse{Type: "error", Message: "You are missing bookID parameter."}
	} else {
		db := db.SetupDBPostgres()

		helper.PrintMessage("Deleting book from DB")

		_, err := db.Exec("DELETE FROM books where bookID = $1", bookID)
		helper.CheckErr(err)

		response = model.JsonResponse{Type: "success", Message: "The book has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all books
func DeleteBooks(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	db := db.SetupDBPostgres()

	helper.PrintMessage("Deleting all books...")

	_, err := db.Exec("DELETE FROM books")
	helper.CheckErr(err)

	helper.PrintMessage("All books have been deleted successfully!")

	var response = model.JsonResponse{Type: "success", Message: "All books have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
