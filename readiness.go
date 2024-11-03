package main

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
)

type ServerHealthResponse struct {
	Status string `json:"status"`
}

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	helpers.RespondWithJSON(w, http.StatusOK, ServerHealthResponse{
		Status: http.StatusText(http.StatusOK),
	})
	fmt.Printf("Healthcheck done with status %s\n", http.StatusText(http.StatusOK))
}
