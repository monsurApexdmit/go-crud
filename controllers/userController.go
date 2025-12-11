package controllers

import (
    "encoding/json"
    "net/http"
    "golang.org/x/crypto/bcrypt"

    "github.com/go-chi/chi/v5"
    "go-crud/database"
    "go-crud/models"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to retrieve users")
        return
    }

    writeJSON(w, http.StatusOK, "Users retrieved successfully", users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var user models.User

    if err := database.DB.First(&user, id).Error; err != nil {
        writeError(w, http.StatusNotFound, "User not found")
        return
    }

    writeJSON(w, http.StatusOK, "User fetched successfully", user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }

    if user.Username == "" || user.Email == "" {
        writeError(w, http.StatusBadRequest, "Username and email are required")
        return
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to process password")
        return
    }
    user.Password = string(hashedPassword)

    if err := database.DB.Create(&user).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }
    user.Password = ""

    writeJSON(w, http.StatusCreated, "User created successfully", user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var user models.User

    if err := database.DB.First(&user, id).Error; err != nil {
        writeError(w, http.StatusNotFound, "User not found")
        return
    }

    var updatedData models.User
    if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }

    user.Username = updatedData.Username
    user.Email = updatedData.Email

    if err := database.DB.Save(&user).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to update user")
        return
    }

    writeJSON(w, http.StatusOK, "User updated successfully", user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to delete user")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
