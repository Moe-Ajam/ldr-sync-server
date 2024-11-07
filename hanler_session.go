package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/google/uuid"
)

type Session struct {
	Users map[string]float64
}

// maps the session ID to usernames
var sessions = make(map[string]*Session)

type SessionCreationResponse struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}

func (cfg apiConfig) handlerCreateSession(w http.ResponseWriter, r *http.Request) {
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	sessionID := uuid.New().String()
	fmt.Printf("Session ID created for the user %s: %s\n", claims.Username, sessionID)
	sessions[sessionID] = &Session{
		Users: make(map[string]float64),
	}
	fmt.Printf("claim username: %s\n", claims.Username)
	sessions[sessionID].Users[claims.Username] = 0.00
	sessionCreationResponse := SessionCreationResponse{
		Username:  claims.Username,
		SessionID: sessionID,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionCreationResponse)
}

type JoinSessionParams struct {
	RequestUser string `json:"request_user"`
	SessionID   string `json:"session_id"`
}

type SessionJoinResponse struct {
	ConnectedUser string `json:"connected_user"`
	SessionID     string `json:"session_id"`
	Message       string `json:"message"`
}

func (cfg apiConfig) handlerJoinSession(w http.ResponseWriter, r *http.Request) {
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	joinSessionParams := JoinSessionParams{}
	err = decoder.Decode(&joinSessionParams)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	_, ok := sessions[joinSessionParams.SessionID]
	if !ok {
		helpers.RespondWithError(w, http.StatusNotFound, "SessionID is invalid")
		return
	}

	sessions[joinSessionParams.SessionID].Users[claims.Username] = 0.00

	sessionJoinResponse := SessionJoinResponse{
		ConnectedUser: joinSessionParams.RequestUser,
		SessionID:     joinSessionParams.SessionID,
		Message:       claims.Username + " is now connected to " + joinSessionParams.RequestUser,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionJoinResponse)
}
