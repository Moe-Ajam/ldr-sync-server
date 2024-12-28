package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Users        map[string]*websocket.Conn
	PlaybackTime map[string]float64
}

type WebSocketMessage struct {
	Action      string  `json:"action"`
	CurrentTime float64 `json:"current_time"`
	UserID      string  `json:"user_id"`
}

// maps the session ID to usernames
var sessions = make(map[string]*Session)

type SessionCreationResponse struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}

func (cfg apiConfig) handlerCreateSession(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		log.Printf("Could not retrieve claims: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	sessionID := uuid.New().String()
	fmt.Printf("Session ID created for the user %s: %s\n", claims.Username, sessionID)
	sessions[sessionID] = &Session{
		Users:        make(map[string]*websocket.Conn),
		PlaybackTime: make(map[string]float64),
	}
	sessions[sessionID].PlaybackTime[claims.Username] = 0.00
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
	SessionID string `json:"session_id"`
}

func (cfg apiConfig) handlerJoinSession(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		log.Printf("Could not retrieve claims: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
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

	sessions[joinSessionParams.SessionID].PlaybackTime[claims.Username] = 0.00

	sessionJoinResponse := SessionJoinResponse{
		SessionID: joinSessionParams.SessionID,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionJoinResponse)
}

func (cfg *apiConfig) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	claims := Claims{}

	sessionID := r.URL.Query().Get("session_id")
	tknString := r.URL.Query().Get("token")

	_, err := jwt.ParseWithClaims(tknString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.jwtSecret), nil
	})
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
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		conn.Close()
		return
	}

	session, exists := sessions[sessionID]
	if !exists {
		helpers.RespondWithError(w, http.StatusNotFound, "Session doesn't exist")
		conn.Close()
		return
	}

	session.Users[claims.Username] = conn
	session.PlaybackTime[claims.Email] = 0.0

	defer func() {
		conn.Close()
		delete(session.Users, claims.Username)
		delete(session.PlaybackTime, claims.Username)

		// Deletes the session if no users are connected to it anymore
		if len(session.Users) == 0 {
			delete(sessions, sessionID)
		}
	}()

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		broadcastToSession(sessionID, msg, claims.Username)
	}
}

func broadcastToSession(sessionID string, msg WebSocketMessage, senderID string) {
	session, exists := sessions[sessionID]
	if !exists {
		return
	}

	for userID, conn := range session.Users {
		if userID != senderID {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending WebSocket message to user %s: %v\n", userID, err)
				conn.Close()
				delete(session.Users, userID)
			}
		}
	}
}
