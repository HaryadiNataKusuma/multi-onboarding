package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// Helper functions for Echo handlers

// RespondWithJSON sends a JSON response
func RespondWithJSON(c echo.Context, code int, payload interface{}) error {
	return c.JSON(code, payload)
}

// RespondWithError sends an error JSON response
func RespondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

// Legacy helper functions for backward compatibility (if needed)
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// This is kept for legacy handlers that might still use http.ResponseWriter
	// New handlers should use RespondWithJSON with echo.Context
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

