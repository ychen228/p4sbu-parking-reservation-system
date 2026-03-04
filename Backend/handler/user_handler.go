package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/middleware"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user related requests
type UserHandler struct {
	repo *mongodb.Repository
}

// NewUserHandler creates a new User handler
func NewUserHandler(repo *mongodb.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetAllUsers handles GET /users (admin only)
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// This endpoint is only for admins, which is enforced in routes.go

	// TODO: Implement get all users in repository
	// For now, return an empty array
	users, err := h.repo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Remove sensitive information
	for i := range users {
		users[i].PasswordHash = ""
	}

	// Return users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get user from repository
	user, err := h.repo.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get the ID of the authenticated user
	authenticatedUserID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user role
	role, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or requesting their own information
	if role != "admin" && authenticatedUserID != id {
		http.Error(w, "Unauthorized to view this user", http.StatusForbidden)
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	// Return user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser handles POST /users (registration)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	existingUser, err := h.repo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Set default role for new users
	if user.Role == "" {
		user.Role = "user"
	}

	// Create new user
	id, err := h.repo.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateUser handles PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get the ID of the authenticated user
	authenticatedUserID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user role
	role, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or updating their own information
	if role != "admin" && authenticatedUserID != id {
		http.Error(w, "Unauthorized to update this user", http.StatusForbidden)
		return
	}

	// Get existing user to preserve some fields
	existingUser, err := h.repo.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	user.ID = id

	// If password is being updated, hash it
	if user.PasswordHash != "" && user.PasswordHash != existingUser.PasswordHash {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		// Keep the existing password
		user.PasswordHash = existingUser.PasswordHash
	}

	// Preserve the role if not an admin
	if role != "admin" {
		user.Role = existingUser.Role
	}

	// Update user
	if err := h.repo.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// This endpoint is only for admins, which is enforced in routes.go

	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete user
	if err := h.repo.DeleteUser(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
