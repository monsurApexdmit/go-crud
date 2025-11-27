package controllers

import (
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go-crud/database"
	"go-crud/models"
)

func ListBooks(w http.ResponseWriter, r *http.Request) {

	rows, err := database.DB.Query("SELECT id, title, author FROM books")
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

func GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var b models.Book
	err = database.DB.QueryRow("SELECT id, title, author FROM books WHERE id = ?", id).Scan(
		&b.ID, &b.Title, &b.Author,
	)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(b)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	json.NewDecoder(r.Body).Decode(&b)

	res, err := database.DB.Exec("INSERT INTO books (title, author) VALUES (?, ?)", b.Title, b.Author)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	b.ID = id

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(b)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var b models.Book
	json.NewDecoder(r.Body).Decode(&b)

	_, err = database.DB.Exec("UPDATE books SET title=?, author=? WHERE id=?",
		b.Title, b.Author, id,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b.ID = id
	json.NewEncoder(w).Encode(b)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	_, err = database.DB.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}
