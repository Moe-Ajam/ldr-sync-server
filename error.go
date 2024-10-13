package main

import "net/http"

func handlerError(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
