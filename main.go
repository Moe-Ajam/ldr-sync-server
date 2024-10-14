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
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not load database: %v\n", err)
	}
	dbQuries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQuries,
	}
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/err", handlerError)

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
