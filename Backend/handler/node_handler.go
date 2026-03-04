package handlers

import (
	"encoding/json"
	"net/http"

	mongodb "github.com/416Coders/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NodeHandler handles node related requests
type NodeHandler struct {
	repo *mongodb.Repository
}

// NewNodeHandler creates a new Node handler
func NewNodeHandler(repo *mongodb.Repository) *NodeHandler {
	return &NodeHandler{repo: repo}
}

// GetNode handles GET /nodes/{id}
func (h *NodeHandler) GetNode(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Get node from repository
	node, err := h.repo.GetNode(id)
	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	// Return node as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

// GetAllNodes handles GET /nodes
func (h *NodeHandler) GetAllNodes(w http.ResponseWriter, r *http.Request) {
	// Get all nodes from repository
	nodes, err := h.repo.GetAllNodes()
	if err != nil {
		http.Error(w, "Failed to retrieve nodes", http.StatusInternalServerError)
		return
	}

	// Return nodes as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

// CreateNode handles POST /nodes
func (h *NodeHandler) CreateNode(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var node models.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new node
	id, err := h.repo.CreateNode(node)
	if err != nil {
		http.Error(w, "Failed to create node", http.StatusInternalServerError)
		return
	}

	// Return new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// UpdateNode handles PUT /nodes/{id}
func (h *NodeHandler) UpdateNode(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Parse request body
	var node models.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	node.ID = id

	// Update node
	if err := h.repo.UpdateNode(node); err != nil {
		http.Error(w, "Failed to update node", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Node updated successfully"})
}

// DeleteNode handles DELETE /nodes/{id}
func (h *NodeHandler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete node
	if err := h.repo.DeleteNode(id); err != nil {
		http.Error(w, "Failed to delete node", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Node deleted successfully"})
}

// GetNavigationPath handles GET /navigate?from={fromNodeId}&to={toNodeId}
func (h *NodeHandler) GetNavigationPath(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	fromID := r.URL.Query().Get("from")
	toID := r.URL.Query().Get("to")

	if fromID == "" || toID == "" {
		http.Error(w, "Both 'from' and 'to' node IDs are required", http.StatusBadRequest)
		return
	}

	// Convert string IDs to ObjectIDs - just validate them for now
	_, err := primitive.ObjectIDFromHex(fromID)
	if err != nil {
		http.Error(w, "Invalid 'from' node ID format", http.StatusBadRequest)
		return
	}

	_, err = primitive.ObjectIDFromHex(toID)
	if err != nil {
		http.Error(w, "Invalid 'to' node ID format", http.StatusBadRequest)
		return
	}

	// TODO: Implement navigation path finding algorithm using fromID and toID
	// For now, return a placeholder response

	// Return an empty path for now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Node{})
}
