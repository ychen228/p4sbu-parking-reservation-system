package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"

	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParkingLotDistance represents a parking lot with its distance from a reference point
type ParkingLotDistance struct {
	ParkingLot models.ParkingLot `json:"parking_lot"`
	Distance   float64           `json:"distance_km"`
}

// WayfindHandler handles wayfinding related requests
type WayfindHandler struct {
	repo *mongodb.Repository
}

// NewWayfindHandler creates a new Wayfind handler
func NewWayfindHandler(repo *mongodb.Repository) *WayfindHandler {
	return &WayfindHandler{repo: repo}
}

// FindNearestParkingLot handles GET /wayfind/building/{id}/nearest-parking
func (h *WayfindHandler) FindNearestParkingLot(w http.ResponseWriter, r *http.Request) {
	// Parse building ID from URL
	vars := mux.Vars(r)
	buildingID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid building ID format", http.StatusBadRequest)
		return
	}

	// Get building from repository
	building, err := h.repo.GetBuilding(buildingID)
	if err != nil {
		http.Error(w, "Building not found", http.StatusNotFound)
		return
	}

	log.Printf("Building found: %s at location lat=%f, lng=%f",
		building.Name, building.Location.Lat, building.Location.Lng)

	// Get all parking lots
	parkingLots, err := h.repo.GetAllParkingLots()
	if err != nil {
		http.Error(w, "Failed to retrieve parking lots", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d total parking lots", len(parkingLots))

	// Calculate distances - include ALL parking lots without filtering
	var parkingLotDistances []ParkingLotDistance
	for _, lot := range parkingLots {
		// Calculate distance for every parking lot
		distance := calculateDistance(building.Location, lot.Location)

		log.Printf("Parking lot: %s, Location: lat=%f, lng=%f, Distance: %f",
			lot.Name, lot.Location.Lat, lot.Location.Lng, distance)

		parkingLotDistances = append(parkingLotDistances, ParkingLotDistance{
			ParkingLot: lot,
			Distance:   distance,
		})
	}

	// Sort by distance
	sort.Slice(parkingLotDistances, func(i, j int) bool {
		return parkingLotDistances[i].Distance < parkingLotDistances[j].Distance
	})

	// Get limit parameter (default to all)
	limit := len(parkingLotDistances)
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 && parsedLimit < limit {
			limit = parsedLimit
		}
	}

	// Return only up to the limit
	result := parkingLotDistances
	if len(result) > limit {
		result = result[:limit]
	}

	log.Printf("Returning %d parking lots sorted by distance", len(result))

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// FindNearestParkingLotByCoordinates handles GET /wayfind/nearest-parking
func (h *WayfindHandler) FindNearestParkingLotByCoordinates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for coordinates
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Missing lat or lng parameters", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid lat parameter", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid lng parameter", http.StatusBadRequest)
		return
	}

	location := models.Location{
		Lat: lat,
		Lng: lng,
	}

	log.Printf("Using coordinates: lat=%f, lng=%f", location.Lat, location.Lng)

	// Get all parking lots
	parkingLots, err := h.repo.GetAllParkingLots()
	if err != nil {
		http.Error(w, "Failed to retrieve parking lots", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d total parking lots", len(parkingLots))

	// Calculate distances - include ALL parking lots without filtering
	var parkingLotDistances []ParkingLotDistance
	for _, lot := range parkingLots {
		// Calculate distance for every parking lot
		distance := calculateDistance(location, lot.Location)

		log.Printf("Parking lot: %s, Location: lat=%f, lng=%f, Distance: %f",
			lot.Name, lot.Location.Lat, lot.Location.Lng, distance)

		parkingLotDistances = append(parkingLotDistances, ParkingLotDistance{
			ParkingLot: lot,
			Distance:   distance,
		})
	}

	// Sort by distance
	sort.Slice(parkingLotDistances, func(i, j int) bool {
		return parkingLotDistances[i].Distance < parkingLotDistances[j].Distance
	})

	// Get limit parameter (default to all)
	limit := len(parkingLotDistances)
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 && parsedLimit < limit {
			limit = parsedLimit
		}
	}

	// Return only up to the limit
	result := parkingLotDistances
	if len(result) > limit {
		result = result[:limit]
	}

	log.Printf("Returning %d parking lots sorted by distance", len(result))

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// calculateDistance calculates the Euclidean distance between two locations
func calculateDistance(loc1, loc2 models.Location) float64 {
	// Using the Euclidean distance formula: sqrt((x2-x1)^2 + (y2-y1)^2)
	dx := loc2.Lng - loc1.Lng
	dy := loc2.Lat - loc1.Lat
	return math.Sqrt(dx*dx + dy*dy)
}
