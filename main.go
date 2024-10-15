package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Moe-Ajam/ldr-sync-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB        *database.Queries
	jwtSecret string
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
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
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
