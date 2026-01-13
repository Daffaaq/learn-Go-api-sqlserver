// ---------------- helpers/helpers.go ----------------
package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string, err interface{}) {
	payload := map[string]interface{}{
		"error": message,
	}
	if err != nil {
		payload["details"] = err
	}
	RespondWithJSON(w, code, payload)
}
