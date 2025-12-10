package controllers

import (
	"encoding/json"
	"net/http"

	"go-crud/database"
	"go-crud/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)


func ListAuthors(w http.ResponseWriter, r *http.Request) {
	var authors []models.Author
	result := database.DB.Find(&authors)
	if result.Error != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve authors")
		return
	}
	writeJSON(w, http.StatusOK, "Authors retrieved successfully", authors)
}


func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author models.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result := database.DB.Create(&author)

	if result.Error != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create author")
		return
	}
	
	writeJSON(w, http.StatusCreated, "Author created successfully", author)
}


func GetAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")


	var author models.Author
	result := database.DB.First(&author, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			writeError(w, http.StatusNotFound, "Author not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to retrieve author")
		}
		return
	}
	writeJSON(w, http.StatusOK, "Author retrieved successfully", author)
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var author models.Author

	if err := database.DB.First(&author, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Author not found")
		return
	}
	
	var updatedData models.Author
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	author.Name = updatedData.Name
	author.Email = updatedData.Email

	if err := database.DB.Save(&author).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update author")
		return
	}

	writeJSON(w, http.StatusOK, "Author updated successfully", author)
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var author models.Author

	if err := database.DB.First(&author, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Author not found")
		return
	}

	if err := database.DB.Delete(&author).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete author")
		return
	}

	writeJSON(w, http.StatusOK, "Author deleted successfully", nil)
}

