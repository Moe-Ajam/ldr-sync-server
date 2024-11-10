package main

import (
	"log"
	"net/http"
)

func enableCORS(w *http.ResponseWriter) {
	log.Println("CORS activated")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
