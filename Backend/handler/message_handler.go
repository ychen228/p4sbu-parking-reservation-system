package handlers

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/YimingChen/P4SBU/Backend/models"
// 	"github.com/YimingChen/P4SBU/Backend/database"
// 	"github.com/gorilla/mux"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type MessageHandler struct {
// 	repo *mongodb.Repository
// }

// func NewMessageHandler(repo *mongodb.Repository) *MessageHandler {
// 	return &MessageHandler{repo: repo}
// }

// func (h *MessageHandler) GeUserMessage(w http.ResponseWriter, r*http.Request) {
// 	vars := mux.Vars(r)
// 	userID, err := primitive.ObjectIDFromHex(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid ID format". http.StatusBadRequest)
// 		return
// 	}

// 	authenticatedUserID, err := middleware.GetUserID(r)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	role, err := middleware.GetUserRole(r)
// 	if err != nil {
// 		http.error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	messages, err := h.repo.GetMessageByUser(userID)
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(messages)
// }

// // func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.request) {
// // 	var message models.Message
// // 	if err := json.NewDecoder(r.body).Decode(&message); err != nil {
// // 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// // 		return
// // 	}

// // 	userID, err := middleware.GetUserID(r)
// // 	if err != nil {
// // 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// // 		return
// // 	}

// // 	message.User = userID

// // }
