package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type createFeedParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (cfg *apiConfig) handler_create_feed(w http.ResponseWriter, r *http.Request, user database.User) {
	reqParams, err := decondeJSON[createFeedParams](r)
	if err != nil {
		fmt.Printf("handler_create_feed decondeJSON error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.Background()
	feed, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      reqParams.Name,
		Url:       reqParams.Url,
		UserID:    user.ID,
	})

	if err != nil {
		fmt.Printf("handler_create_feed CreateFeed error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, feed)
}
