package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/416Coders/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParkingLotHandler handles parking lot related requests
type ParkingLotHandler struct {
	repo *mongodb.Repository
}

// NewParkingLotHandler creates a new ParkingLot handler
func NewParkingLotHandler(repo *mongodb.Repository) *ParkingLotHandler {
	return &ParkingLotHandler{repo: repo}
}

// GetParkingLot handles GET /parking-lots/{id}
func (h *ParkingLotHandler) GetParkingLot(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get parking lot from repository
	parkingLot, err := h.repo.GetParkingLot(id)
	if err != nil {
		http.Error(w, "Parking lot not found", http.StatusNotFound)
		return
	}

	// Return parking lot as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingLot)
}

// GetAllParkingLots handles GET /parking-lots
func (h *ParkingLotHandler) GetAllParkingLots(w http.ResponseWriter, r *http.Request) {
	// Get all parking lots from repository
	parkingLots, err := h.repo.GetAllParkingLots()
	if err != nil {
		http.Error(w, "Failed to retrieve parking lots", http.StatusInternalServerError)
		return
	}

	// Return parking lots as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingLots)
}

// CreateParkingLot handles POST /parking-lots
func (h *ParkingLotHandler) CreateParkingLot(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var parkingLot models.ParkingLot
	if err := json.NewDecoder(r.Body).Decode(&parkingLot); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new parking lot
	id, err := h.repo.CreateParkingLot(parkingLot)
	if err != nil {
		http.Error(w, "Failed to create parking lot", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateParkingLot handles PUT /parking-lots/{id}
func (h *ParkingLotHandler) UpdateParkingLot(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Parse request body
	var parkingLot models.ParkingLot
	if err := json.NewDecoder(r.Body).Decode(&parkingLot); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	parkingLot.ID = id

	// Update parking lot
	if err := h.repo.UpdateParkingLot(parkingLot); err != nil {
		http.Error(w, "Failed to update parking lot", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Parking lot updated successfully"})
}

// DeleteParkingLot handles DELETE /parking-lots/{id}
func (h *ParkingLotHandler) DeleteParkingLot(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete parking lot
	if err := h.repo.DeleteParkingLot(id); err != nil {
		http.Error(w, "Failed to delete parking lot", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Parking lot deleted successfully"})
}
