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
	_ "github.com/lib/pq"
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
	// port := os.Getenv("PORT")
	port := "8080"
	dbURL := os.Getenv("CONN")
	secret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not load database: %v\n", err)
	}
	dbQuries := database.New(db)

	apiCfg := apiConfig{
		DB:        dbQuries,
		jwtSecret: secret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/err", handlerError)

	mux.HandleFunc("POST /api/register", apiCfg.handlerUserCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("/api/ws", apiCfg.handleConnection)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/create-session", apiCfg.handlerCreateSession)
	mux.HandleFunc("POST /api/join-session", apiCfg.handlerJoinSession)

	go handleMessage()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
