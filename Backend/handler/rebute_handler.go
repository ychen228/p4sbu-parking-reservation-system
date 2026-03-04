package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/416Coders/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/middleware"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ViolationHandler handles violation related requests
type ViolationHandler struct {
	repo *mongodb.Repository
}

// NewViolationHandler creates a new Violation handler
func NewViolationHandler(repo *mongodb.Repository) *ViolationHandler {
	return &ViolationHandler{repo: repo}
}

// GetAllViolations handles GET /violations (admin only)
func (h *ViolationHandler) GetAllViolations(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// TODO: Implement get all violations in repository
	// For now, returning empty array
	violations := []models.Violation{}

	// Return violations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violations)
}

// GetViolation handles GET /violations/{id}
func (h *ViolationHandler) GetViolation(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get violation from repository
	violation, err := h.repo.GetViolation(id)
	if err != nil {
		http.Error(w, "Violation not found", http.StatusNotFound)
		return
	}

	// Get the ID of the authenticated user
	userID, err := middleware.GetUserID(r)
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

	// Check if user is admin or the user associated with the violation
	if role != "admin" && userID.Hex() != violation.User.Hex() {
		http.Error(w, "Unauthorized to view this violation", http.StatusForbidden)
		return
	}

	// Return violation as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violation)
}

// GetUserViolations handles GET /users/{id}/violations
func (h *ViolationHandler) GetUserViolations(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
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

	// Check if user is admin or requesting their own violations
	if role != "admin" && authenticatedUserID != userID {
		http.Error(w, "Unauthorized to view this user's violations", http.StatusForbidden)
		return
	}

	// Get violations for user
	violations, err := h.repo.GetViolationsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve violations", http.StatusInternalServerError)
		return
	}

	// Return violations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violations)
}

// CreateViolation handles POST /admin/violations (admin only)
func (h *ViolationHandler) CreateViolation(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// Parse request body
	var violation models.Violation
	if err := json.NewDecoder(r.Body).Decode(&violation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user exists
	_, err := h.repo.GetUser(violation.User)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Validate the parking lot exists
	_, err = h.repo.GetParkingLot(violation.ParkingLot)
	if err != nil {
		http.Error(w, "Parking lot not found", http.StatusBadRequest)
		return
	}

	// Create new violation
	id, err := h.repo.CreateViolation(violation)
	if err != nil {
		http.Error(w, "Failed to create violation", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateViolation handles PUT /admin/violations/{id} (admin only)
func (h *ViolationHandler) UpdateViolation(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check if violation exists
	_, err = h.repo.GetViolation(id)
	if err != nil {
		http.Error(w, "Violation not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var violation models.Violation
	if err := json.NewDecoder(r.Body).Decode(&violation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	violation.ID = id

	// Update violation
	if err := h.repo.UpdateViolation(violation); err != nil {
		http.Error(w, "Failed to update violation", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Violation updated successfully"})
}

// DeleteViolation handles DELETE /admin/violations/{id} (admin only)
func (h *ViolationHandler) DeleteViolation(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete violation
	if err := h.repo.DeleteViolation(id); err != nil {
		http.Error(w, "Failed to delete violation", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Violation deleted successfully"})
}

// ViolationRebuteHandler handles violation rebute related requests
type ViolationRebuteHandler struct {
	repo *mongodb.Repository
}

// NewViolationRebuteHandler creates a new ViolationRebute handler
func NewViolationRebuteHandler(repo *mongodb.Repository) *ViolationRebuteHandler {
	return &ViolationRebuteHandler{repo: repo}
}

// GetAllRebutes handles GET /violation-rebutes (admin only)
func (h *ViolationRebuteHandler) GetAllRebutes(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// TODO: Implement get all rebutes in repository
	// For now, returning empty array
	rebutes := []models.ViolationRebute{}

	// Return rebutes as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rebutes)
}

// GetRebute handles GET /violation-rebutes/{id}
func (h *ViolationRebuteHandler) GetRebute(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get rebute from repository
	rebute, err := h.repo.GetViolationRebute(id)
	if err != nil {
		http.Error(w, "Rebute not found", http.StatusNotFound)
		return
	}

	// Get the ID of the authenticated user
	userID, err := middleware.GetUserID(r)
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

	// Check if user is admin or the user associated with the rebute
	if role != "admin" && userID.Hex() != rebute.User.Hex() {
		http.Error(w, "Unauthorized to view this rebute", http.StatusForbidden)
		return
	}

	// Return rebute as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rebute)
}

// GetUserRebutes handles GET /users/{id}/violation-rebutes
func (h *ViolationRebuteHandler) GetUserRebutes(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["id"])
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

	// Check if user is admin or requesting their own rebutes
	if role != "admin" && authenticatedUserID != userID {
		http.Error(w, "Unauthorized to view this user's rebutes", http.StatusForbidden)
		return
	}

	// Get rebutes for user
	rebutes, err := h.repo.GetViolationRebutesByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve rebutes", http.StatusInternalServerError)
		return
	}

	// Return rebutes as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rebutes)
}

// GetViolationRebutes handles GET /violations/{id}/rebutes
func (h *ViolationRebuteHandler) GetViolationRebutes(w http.ResponseWriter, r *http.Request) {
	// Parse violation ID from URL
	vars := mux.Vars(r)
	violationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get the violation to check permissions
	violation, err := h.repo.GetViolation(violationID)
	if err != nil {
		http.Error(w, "Violation not found", http.StatusNotFound)
		return
	}

	// Get the ID of the authenticated user
	userID, err := middleware.GetUserID(r)
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

	// Check if user is admin or the user associated with the violation
	if role != "admin" && userID.Hex() != violation.User.Hex() {
		http.Error(w, "Unauthorized to view rebutes for this violation", http.StatusForbidden)
		return
	}

	// TODO: Implement get rebutes by violation ID in repository
	// For now, returning empty array
	rebutes := []models.ViolationRebute{}

	// Return rebutes as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rebutes)
}

// CreateRebute handles POST /violation-rebutes
func (h *ViolationRebuteHandler) CreateRebute(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var rebute models.ViolationRebute
	if err := json.NewDecoder(r.Body).Decode(&rebute); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the ID of the authenticated user
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set the user ID as the rebute submitter
	rebute.User = userID

	// Validate the violation exists
	violation, err := h.repo.GetViolation(rebute.Violation)
	if err != nil {
		http.Error(w, "Violation not found", http.StatusBadRequest)
		return
	}

	// Check that the user is the one associated with the violation
	if userID.Hex() != violation.User.Hex() {
		http.Error(w, "Cannot submit rebute for another user's violation", http.StatusForbidden)
		return
	}

	// Set the parking lot from the violation
	rebute.ParkingLot = violation.ParkingLot

	// Create new rebute
	id, err := h.repo.CreateViolationRebute(rebute)
	if err != nil {
		http.Error(w, "Failed to create rebute", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateRebute handles PUT /violation-rebutes/{id} (admin only)
func (h *ViolationRebuteHandler) UpdateRebute(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check if rebute exists
	existingRebute, err := h.repo.GetViolationRebute(id)
	if err != nil {
		http.Error(w, "Rebute not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var rebute models.ViolationRebute
	if err := json.NewDecoder(r.Body).Decode(&rebute); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	rebute.ID = id

	// Preserve the original user, parking lot, and violation
	rebute.User = existingRebute.User
	rebute.ParkingLot = existingRebute.ParkingLot
	rebute.Violation = existingRebute.Violation

	// Update rebute
	if err := h.repo.UpdateViolationRebute(rebute); err != nil {
		http.Error(w, "Failed to update rebute", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rebute updated successfully"})
}

// DeleteRebute handles DELETE /violation-rebutes/{id} (admin only)
func (h *ViolationRebuteHandler) DeleteRebute(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete rebute
	if err := h.repo.DeleteViolationRebute(id); err != nil {
		http.Error(w, "Failed to delete rebute", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rebute deleted successfully"})
}
