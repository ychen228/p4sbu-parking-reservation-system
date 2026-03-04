package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/middleware"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReservationHandler handles reservation related requests
type ReservationHandler struct {
	repo *mongodb.Repository
}

// NewReservationHandler creates a new Reservation handler
func NewReservationHandler(repo *mongodb.Repository) *ReservationHandler {
	return &ReservationHandler{repo: repo}
}

// GetAllReservations handles GET /reservations (admin only)
func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	// This is already protected by admin middleware in routes.go

	// TODO: Implement get all reservations in repository
	// For now, return an empty array
	reservations, err := h.repo.GetAllReservations()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Return reservations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

// GetReservation handles GET /reservations/{id}
func (h *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get reservation from repository
	reservation, err := h.repo.GetReservation(id)
	if err != nil {
		http.Error(w, "Reservation not found", http.StatusNotFound)
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

	// Check if user is admin or the owner of the reservation
	if role != "admin" && userID.Hex() != reservation.ReservedBy.Hex() {
		http.Error(w, "Unauthorized to view this reservation", http.StatusForbidden)
		return
	}

	// Return reservation as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

// GetUserReservations handles GET /users/{id}/reservations
func (h *ReservationHandler) GetUserReservations(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is admin or requesting their own reservations
	if role != "admin" && authenticatedUserID != userID {
		http.Error(w, "Unauthorized to view this user's reservations", http.StatusForbidden)
		return
	}

	// Get reservations for user
	reservations, err := h.repo.GetReservationsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve reservations", http.StatusInternalServerError)
		return
	}

	// Return reservations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

// GetParkingLotReservations handles GET /parking-lots/{id}/reservations
func (h *ReservationHandler) GetParkingLotReservations(w http.ResponseWriter, r *http.Request) {
	// Parse parking lot ID from URL
	vars := mux.Vars(r)
	parkingLotID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// This endpoint is only for admins, which is enforced in routes.go

	// Get reservations for parking lot
	reservations, err := h.repo.GetReservationsByLot(parkingLotID)
	if err != nil {
		http.Error(w, "Failed to retrieve reservations", http.StatusInternalServerError)
		return
	}

	// Return reservations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

// CreateReservation handles POST /reservations
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the ID of the authenticated user
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set the user ID as the reservation owner
	reservation.ReservedBy = userID

	// Validate reservation times
	if reservation.StartTime.Before(time.Now()) {
		http.Error(w, "Start time must be in the future", http.StatusBadRequest)
		return
	}

	if reservation.EndTime.Before(reservation.StartTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Check if the parking lot exists
	_, err = h.repo.GetParkingLot(reservation.ParkingLot)
	if err != nil {
		http.Error(w, "Parking lot not found", http.StatusBadRequest)
		return
	}

	// TODO: Check for conflicting reservations

	// Set initial status
	reservation.Status = "confirmed"

	// Create new reservation
	id, err := h.repo.CreateReservation(reservation)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateReservation handles PUT /reservations/{id}
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get existing reservation
	existingReservation, err := h.repo.GetReservation(id)
	if err != nil {
		http.Error(w, "Reservation not found", http.StatusNotFound)
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

	// Check if user is admin or the owner of the reservation
	if role != "admin" && userID.Hex() != existingReservation.ReservedBy.Hex() {
		http.Error(w, "Unauthorized to update this reservation", http.StatusForbidden)
		return
	}

	// Parse request body
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	reservation.ID = id

	// Preserve the original owner
	reservation.ReservedBy = existingReservation.ReservedBy

	// Validate reservation times
	if reservation.StartTime.Before(time.Now()) {
		http.Error(w, "Start time must be in the future", http.StatusBadRequest)
		return
	}

	if reservation.EndTime.Before(reservation.StartTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Update reservation
	if err := h.repo.UpdateReservation(reservation); err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully"})
}

// CancelReservation handles DELETE /reservations/{id}
func (h *ReservationHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get existing reservation
	existingReservation, err := h.repo.GetReservation(id)
	if err != nil {
		http.Error(w, "Reservation not found", http.StatusNotFound)
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

	// Check if user is admin or the owner of the reservation
	if role != "admin" && userID.Hex() != existingReservation.ReservedBy.Hex() {
		http.Error(w, "Unauthorized to cancel this reservation", http.StatusForbidden)
		return
	}

	// Delete reservation
	if err := h.repo.DeleteReservation(id); err != nil {
		http.Error(w, "Failed to cancel reservation", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation cancelled successfully"})
}
