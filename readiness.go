package main

import (
	"fmt"
	"net/http"
)

type ServerHealthResponse struct {
	Status string `json:"status"`
}

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, ServerHealthResponse{
		Status: http.StatusText(http.StatusOK),
	})
	fmt.Printf("Healthcheck done with status %s\n", http.StatusText(http.StatusOK))
}
