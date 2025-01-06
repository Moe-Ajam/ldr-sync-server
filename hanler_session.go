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

// Maps session IDs to session data
var sessions = make(map[string]*Session)

type SessionCreationResponse struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}

type JoinSessionParams struct {
	SessionID string `json:"session_id"`
}

type SessionJoinResponse struct {
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

	// Create a new session with a unique ID
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

	// Check if the session exists
	session, ok := sessions[joinSessionParams.SessionID]
	if !ok {
		helpers.RespondWithError(w, http.StatusNotFound, "SessionID is invalid")
		return
	}

	// Add the user to the session's playback map
	session.PlaybackTime[claims.Username] = 0.00

	sessionJoinResponse := SessionJoinResponse{
		SessionID: joinSessionParams.SessionID,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, &sessionJoinResponse)
}

func (cfg *apiConfig) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to connect to socket...")
	enableCORS(&w, r)
	claims := Claims{}

	// Extract session ID and token from the query parameters
	sessionID := r.URL.Query().Get("session_id")
	tknString := r.URL.Query().Get("token")
	log.Printf("Connecting to web socket with session ID: %s, and the token received is: %s\n", sessionID, tknString)

	// Parse the JWT token
	_, err := jwt.ParseWithClaims(tknString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Unauthorized - Invalid token signature")
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}
		log.Println("Unauthorized - Token parsing error")
		http.Error(w, "Unauthorized - Token parsing error", http.StatusUnauthorized)
		return
	}

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Something went wrong while upgrading the connection: %v\n", err)
		conn.Close()
		return
	}
	log.Println("Connection upgraded successfully!")

	// Retrieve the session
	session, exists := sessions[sessionID]
	if !exists {
		log.Println("Session does not exist")
		conn.Close()
		return
	}

	// Add the user's connection to the session
	session.Users[claims.Username] = conn
	session.PlaybackTime[claims.Username] = 0.0

	// Ensure cleanup on disconnect
	defer func() {
		conn.Close()
		delete(session.Users, claims.Username)
		delete(session.PlaybackTime, claims.Username)

		// Delete the session if no users are connected
		if len(session.Users) == 0 {
			delete(sessions, sessionID)
		}
	}()

	// Listen for messages from the client
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		log.Printf("Message received from %s: %v\n", claims.Username, msg)
		broadcastToSession(sessionID, msg, claims.Username)
	}
}

func broadcastToSession(sessionID string, msg WebSocketMessage, senderID string) {
	session, exists := sessions[sessionID]
	if !exists {
		log.Printf("Session %s does not exist\n", sessionID)
		return
	}

	for userID, conn := range session.Users {
		if userID != senderID {
			log.Printf("Sending message from %s to %s: %v\n", senderID, userID, msg)
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending WebSocket message to user %s: %v\n", userID, err)
				conn.Close()
				delete(session.Users, userID)
			}
		}
	}
}
