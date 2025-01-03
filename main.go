package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Moe-Ajam/ldr-sync-server/internal/database"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	// _ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	DB        *database.Queries
	jwtSecret string
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	// port := "8080"
	// dbURL := os.Getenv("CONN")
	secret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("could not load database: %v\n", err)
	}
	dbQuries := database.New(db)

	apiCfg := apiConfig{
		DB:        dbQuries,
		jwtSecret: secret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/err", handlerError)

	mux.HandleFunc("/api/register", apiCfg.handlerUserCreate)
	mux.HandleFunc("/api/login", apiCfg.handlerUserLogin)

	mux.HandleFunc("/api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("/api/create-session", apiCfg.handlerCreateSession)
	mux.HandleFunc("/api/join-session", apiCfg.handlerJoinSession)
	mux.HandleFunc("/api/ws", apiCfg.handlerWebSocket)
	mux.HandleFunc("/api/validate-token", apiCfg.hanlderValidateToken)

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}
	log.Println("Starting server on https://api.moecodes.com...")
	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServeTLS("/etc/letsencrypt/live/api.moecodes.com/fullchain.pem", "/etc/letsencrypt/live/api.moecodes.com/privkey.pem"))
	// log.Fatal(srv.ListenAndServe())
}
