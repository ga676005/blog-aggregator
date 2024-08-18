package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	dbURL := os.Getenv("dbURL")
	if dbURL == "" {
		log.Fatal("dbURL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't connect to db %v", err)
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	feedWorker := NewFeedWorker(time.Minute, 10, &cfg)
	go feedWorker.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		text := "Hello world!"
		w.Write([]byte(text))
	})

	mux.HandleFunc("/v1/healthz", handler_readiness)
	mux.HandleFunc("/v1/err", handler_error)
	mux.HandleFunc("POST /v1/users", cfg.handler_create_user)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handler_get_user))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handler_create_feed))
	mux.HandleFunc("GET /v1/feeds", cfg.getAllFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handlerPostUserFollowFeed))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handlerDeleteUserFollowFeed))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.handlerGetUserFollowFeed))

	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.handlerGetPostsByUser))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("starting server at %s", server.Addr)
	server.ListenAndServe()
}
