package controllers

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "go-crud/database"
    "go-crud/models"
)

func ListBooks(w http.ResponseWriter, r *http.Request) {
    var books []models.Book
    if err := database.DB.Find(&books).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to retrieve books")
        return
    }

    writeJSON(w, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        writeError(w, http.StatusNotFound, "Book not found")
        return
    }

    writeJSON(w, http.StatusOK, "Book fetched successfully", book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book

    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }

    if err := database.DB.Create(&book).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to create book")
        return
    }

    writeJSON(w, http.StatusCreated, "Book created successfully", book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var book models.Book

    if err := database.DB.First(&book, id).Error; err != nil {
        writeError(w, http.StatusNotFound, "Book not found")
        return
    }

    var updatedData models.Book
    if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }

    // Update fields manually (recommended)
    book.Title = updatedData.Title
    book.Author = updatedData.Author

    if err := database.DB.Save(&book).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to update book")
        return
    }

    writeJSON(w, http.StatusOK, "Book updated successfully", book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := database.DB.Delete(&models.Book{}, id).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to delete book")
        return
    }

    // No body should be sent for 204 response
    w.WriteHeader(http.StatusNoContent)
}

//
// Helper functions
//

func writeJSON(w http.ResponseWriter, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    json.NewEncoder(w).Encode(models.Response{
        Status:  status,
        Message: message,
        Data:    data,
    })
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, msg, nil)
}
