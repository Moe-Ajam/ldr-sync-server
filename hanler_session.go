package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// maps the session ID to usernames
var sessions = make(map[string][]string)

type SessionCreationResponse struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}

func (cfg apiConfig) handlerCreateSession(w http.ResponseWriter, r *http.Request) {
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return
		}
		if err == jwt.ErrSignatureInvalid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token is invalid")
		fmt.Println("Something went wrong!!")
		return
	}

	sessionID := uuid.New().String()
	fmt.Printf("Session ID created for the user %s: %s\n", claims.Username, sessionID)
	sessions[sessionID] = []string{}
	sessions[sessionID] = append(sessions[sessionID], claims.Username)
	sessionCreationResponse := SessionCreationResponse{
		Username:  claims.Username,
		SessionID: sessionID,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionCreationResponse)
}

type JoinSessionParams struct {
	SessionID string `json:"session_id"`
}

type SessionJoinResponse struct {
	ConnectedUser string `json:"connected_user"`
	SessionID     string `json:"session_id"`
}

func (cfg apiConfig) handlerJoinSession(w http.ResponseWriter, r *http.Request) {
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return
		}
		if err == jwt.ErrSignatureInvalid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token is invalid")
		fmt.Println("Something went wrong!!")
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

	session, ok := sessions[joinSessionParams.SessionID]
	if !ok {
		helpers.RespondWithError(w, http.StatusNotFound, "SessionID is invalid")
		return
	}

	sessions[joinSessionParams.SessionID] = append(sessions[joinSessionParams.SessionID], claims.Username)
	sessionJoinResponse := SessionJoinResponse{
		ConnectedUser: session[0],
		SessionID:     joinSessionParams.SessionID,
	}
	helpers.RespondWithJSON(w, http.StatusCreated, &sessionJoinResponse)
}
