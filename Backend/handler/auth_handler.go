package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	repo *mongodb.Repository
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	UserID string      `json:"userId"`
	Role   string      `json:"role"`
	User   models.User `json:"user"`
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(repo *mongodb.Repository) *AuthHandler {
	return &AuthHandler{
		repo: repo,
	}
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by username
	user, err := h.repo.GetUserByUsername(loginReq.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""

	// Create response
	response := LoginResponse{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		User:   *user,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	_, err := h.repo.GetUserByUsername(user.Username)
	if err == nil {
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

	// Set default role if not specified
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
