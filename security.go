package main

import (
	"net/http"
)

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "chrome-extension://mkjhflenhpjedegkhgnjlconogccecmp")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func (cfg *apiConfig) hanlderValidateToken(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
}
