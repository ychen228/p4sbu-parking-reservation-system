package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/416Coders/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BuildingHandler handles building related requests
type BuildingHandler struct {
	repo *mongodb.Repository
}

// NewBuildingHandler creates a new Building handler
func NewBuildingHandler(repo *mongodb.Repository) *BuildingHandler {
	return &BuildingHandler{repo: repo}
}

// GetBuilding handles GET /buildings/{id}
func (h *BuildingHandler) GetBuilding(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get building from repository
	building, err := h.repo.GetBuilding(id)
	if err != nil {
		http.Error(w, "Building not found", http.StatusNotFound)
		return
	}

	// Return building as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(building)
}

// GetAllBuildings handles GET /buildings
func (h *BuildingHandler) GetAllBuildings(w http.ResponseWriter, r *http.Request) {
	// Get all buildings from repository
	buildings, err := h.repo.GetAllBuildings()
	if err != nil {
		http.Error(w, "Failed to retrieve buildings", http.StatusInternalServerError)
		return
	}

	// Return buildings as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}

// CreateBuilding handles POST /buildings
func (h *BuildingHandler) CreateBuilding(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var building models.Building
	if err := json.NewDecoder(r.Body).Decode(&building); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new building
	id, err := h.repo.CreateBuilding(building)
	if err != nil {
		http.Error(w, "Failed to create building", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateBuilding handles PUT /buildings/{id}
func (h *BuildingHandler) UpdateBuilding(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Parse request body
	var building models.Building
	if err := json.NewDecoder(r.Body).Decode(&building); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	building.ID = id

	// Update building
	if err := h.repo.UpdateBuilding(building); err != nil {
		http.Error(w, "Failed to update building", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Building updated successfully"})
}

// DeleteBuilding handles DELETE /buildings/{id}
func (h *BuildingHandler) DeleteBuilding(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete building
	if err := h.repo.DeleteBuilding(id); err != nil {
		http.Error(w, "Failed to delete building", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Building deleted successfully"})
}
