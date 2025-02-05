package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, `{"error": "Invalid receipt ID format"}`, http.StatusBadRequest)
		return
	}

	receipt, exists := receiptStore[id]
	if !exists {
		http.Error(w, `{"error": "No receipt found for that ID."}`, http.StatusNotFound)
		return
	}

	response := map[string]int{"points": receipt.Points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
