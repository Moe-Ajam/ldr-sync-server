package main

import (
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
)

func handlerError(w http.ResponseWriter, _ *http.Request) {
	helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
