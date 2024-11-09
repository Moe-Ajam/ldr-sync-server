package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
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
	claims := Claims{}
	err := auth.GetClaims(w, r, &claims, cfg.jwtSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	sessionID := uuid.New().String()
	fmt.Printf("Session ID created for the user %s: %s\n", claims.Username, sessionID)
	sessions[sessionID] = &Session{
		Users:        make(map[string]*websocket.Conn),
		PlaybackTime: make(map[string]float64),
	}
	fmt.Printf("claim username: %s\n", claims.Username)
	sessions[sessionID].PlaybackTime[claims.Username] = 0.00
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

	sessions[joinSessionParams.SessionID].PlaybackTime[claims.Username] = 0.00

	sessionJoinResponse := SessionJoinResponse{
		ConnectedUser: joinSessionParams.RequestUser,
		SessionID:     joinSessionParams.SessionID,
		Message:       claims.Username + " is now connected to " + joinSessionParams.RequestUser,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionJoinResponse)
}

func (cfg *apiConfig) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	userID := r.URL.Query().Get("user_id")

	log.Printf("The session ID is: %s and the user is: %s\n", sessionID, userID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		conn.Close()
		return
	}

	session, exists := sessions[sessionID]
	log.Printf("Session is: %v, and it %v", session, exists)
	if !exists {
		helpers.RespondWithError(w, http.StatusNotFound, "Session doesn't exist")
		conn.Close()
		return
	}

	session.Users[userID] = conn
	session.PlaybackTime[userID] = 0.0

	defer func() {
		conn.Close()
		delete(session.Users, userID)
		delete(session.PlaybackTime, userID)
	}()

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		broadcastToSession(sessionID, msg, userID)
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
				log.Printf("Error sending WebSocket messgae to user %s: %v\n", userID, err)
				conn.Close()
				delete(session.Users, userID)
			}
		}
	}
}
