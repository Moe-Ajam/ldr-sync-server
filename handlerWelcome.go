package main

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
)

func (cfg *apiConfig) handleWelcome(w http.ResponseWriter, r *http.Request) {
	claims := &Claims{}
	err := auth.GetClaims(w, r, claims, cfg.jwtSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s, with email %s!", claims.Username, claims.Email)))
}
