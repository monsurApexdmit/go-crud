package controllers

import (
    "encoding/json"
    "net/http"
    // "strconv"

    "github.com/go-chi/chi/v5"
    "go-crud/database"
    "go-crud/models"
)

func ListBooks(w http.ResponseWriter, r *http.Request) {
    var books []models.Book
    database.DB.Find(&books)
    json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        http.Error(w, "Book not found", 404)
        return
    }

    json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    json.NewDecoder(r.Body).Decode(&book)

    database.DB.Create(&book)

    w.WriteHeader(201)
    json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        http.Error(w, "Book not found", 404)
        return
    }

    json.NewDecoder(r.Body).Decode(&book)
    database.DB.Save(&book)

    json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := database.DB.Delete(&models.Book{}, id).Error; err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.WriteHeader(204)
}
