package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/utils"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

var receiptStore = make(map[string]models.Receipt)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var rawRequest map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&rawRequest)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	requiredFields := []string{"retailer", "purchaseDate", "purchaseTime", "items", "total"}
	for _, field := range requiredFields {
		value, exists := rawRequest[field]
		if !exists {
			http.Error(w, fmt.Sprintf(`{"error": "Missing required field: %s"}`, field), http.StatusBadRequest)
			return
		}

		if str, ok := value.(string); ok && str == "" {
			http.Error(w, fmt.Sprintf(`{"error": "Field %s cannot be empty"}`, field), http.StatusBadRequest)
			return
		}
	}

	items, exists := rawRequest["items"].([]interface{})
	if !exists || len(items) == 0 {
		http.Error(w, `{"error": "Items must be a non-empty array"}`, http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseFloat(rawRequest["total"].(string), 64); err != nil {
		http.Error(w, `{"error": "Invalid total format, expected decimal XX.XX"}`, http.StatusBadRequest)
		return
	}

	var receipt models.Receipt
	err = json.Unmarshal([]byte(toJSON(rawRequest)), &receipt)
	if err != nil {
		http.Error(w, `{"error": "Error processing receipt"}`, http.StatusInternalServerError)
		return
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-&]+$`)
	if !re.MatchString(receipt.Retailer) {
		http.Error(w, `{"error": "Retailer name contains invalid characters"}`, http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	points := utils.CalculatePoints(receipt)

	receipt.Points = points
	receiptStore[receiptID] = receipt

	response := map[string]string{"id": receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func toJSON(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
