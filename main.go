package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go-crud/models"
)


var db *sql.DB

func main() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database unreachable: ", err)
	}

	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", bookHandler)

	log.Println("API running on :8004")
	http.ListenAndServe(":8004", nil)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listBooks(w)
	case "POST":
		createBook(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	switch r.Method {
	case "GET":
		getBook(w, id)
	case "PUT":
		updateBook(w, r, id)
	case "DELETE":
		deleteBook(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func listBooks(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author)
		books = append(books, b)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, id int64) {
	var b models.Book
	err := db.QueryRow("SELECT id, title, author FROM books WHERE id = ?", id).Scan(
		&b.ID, &b.Title, &b.Author,
	)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(b)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	json.NewDecoder(r.Body).Decode(&b)

	res, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", b.Title, b.Author)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	b.ID = id

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(b)
}

func updateBook(w http.ResponseWriter, r *http.Request, id int64) {
	var b models.Book
	json.NewDecoder(r.Body).Decode(&b)

	_, err := db.Exec("UPDATE books SET title=?, author=? WHERE id=?",
		b.Title, b.Author, id,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b.ID = id
	json.NewEncoder(w).Encode(b)
}

func deleteBook(w http.ResponseWriter, id int64) {
	_, err := db.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}
