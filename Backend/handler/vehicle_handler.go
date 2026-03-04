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

// VehicleHandler handles vehicle related requests
type VehicleHandler struct {
	repo *mongodb.Repository
}

// NewVehicleHandler creates a new Vehicle handler
func NewVehicleHandler(repo *mongodb.Repository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

// GetVehicle handles GET /vehicles/{id}
func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get vehicle from repository
	vehicle, err := h.repo.GetVehicle(id)
	if err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	// Get the user role for authorization
	role, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or the owner of the vehicle
	// This would need vehicle owner check - for now we can check if admin
	if role != "admin" {
		// Here you would check if the vehicle belongs to the user
		// For now, just checking admin role
		http.Error(w, "Unauthorized to view this vehicle", http.StatusForbidden)
		return
	}

	// Return vehicle as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

// GetUserVehicles handles GET /users/{id}/vehicles
func (h *VehicleHandler) GetUserVehicles(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is admin or requesting their own vehicles
	if role != "admin" && authenticatedUserID != userID {
		http.Error(w, "Unauthorized to view this user's vehicles", http.StatusForbidden)
		return
	}

	// Get user from repository to ensure they exist
	_, err = h.repo.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// TODO: Implement get vehicles by user ID in repository
	// For now, return an empty array
	vehicles := []models.Vehicle{}

	// Return vehicles as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

// CreateVehicle handles POST /vehicles
func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new vehicle
	id, err := h.repo.CreateVehicle(vehicle)
	if err != nil {
		http.Error(w, "Failed to create vehicle", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateVehicle handles PUT /vehicles/{id}
func (h *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get the role for authorization
	role, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or the owner of the vehicle
	if role != "admin" {
		// Here you would check if the vehicle belongs to the user
		// For now, just check if admin
		http.Error(w, "Unauthorized to update this vehicle", http.StatusForbidden)
		return
	}

	// Parse request body
	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	vehicle.ID = id

	// Update vehicle
	if err := h.repo.UpdateVehicle(vehicle); err != nil {
		http.Error(w, "Failed to update vehicle", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle updated successfully"})
}

// DeleteVehicle handles DELETE /vehicles/{id}
func (h *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get the role for authorization
	role, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or the owner of the vehicle
	if role != "admin" {
		// Here you would check if the vehicle belongs to the user
		// For now, just check if admin
		http.Error(w, "Unauthorized to delete this vehicle", http.StatusForbidden)
		return
	}

	// Delete vehicle
	if err := h.repo.DeleteVehicle(id); err != nil {
		http.Error(w, "Failed to delete vehicle", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle deleted successfully"})
}
