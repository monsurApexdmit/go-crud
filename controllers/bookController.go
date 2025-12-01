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

    response := models.Response{
        Status:  http.StatusOK,
        Message: "Books retrieved successfully",
        Data:    books,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        response := models.Response{
            Status:  404,
            Message: "Book not found",
            Data:    nil,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := models.Response{
        Status:  200,
        Message: "Book fetched successfully",
        Data:    book,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


func CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    json.NewDecoder(r.Body).Decode(&book)

    database.DB.Create(&book)

    response := models.Response{
        Status:  201,
        Message: "Book created successfully",
        Data:    book,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}


func UpdateBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        response := models.Response{
            Status:  404,
            Message: "Book not found",
            Data:    nil,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(response)
        return
    }

    json.NewDecoder(r.Body).Decode(&book)
    database.DB.Save(&book)

    response := models.Response{
        Status:  200,
        Message: "Book updated successfully",
        Data:    book,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


func DeleteBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := database.DB.Delete(&models.Book{}, id).Error; err != nil {
        response := models.Response{
            Status:  500,
            Message: "Failed to delete book",
            Data:    nil,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := models.Response{
        Status:  204,
        Message: "Book deleted successfully",
        Data:    nil,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(response)
}

