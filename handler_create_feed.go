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

type CreateFeedResponse struct {
	Feed        database.Feed      `json:"feed"`
	Feed_follow database.UsersFeed `json:"feed_follow"`
}

func (cfg *apiConfig) handler_create_feed(w http.ResponseWriter, r *http.Request, user database.User) {
	reqParams, err := parseRequestJSON[createFeedParams](r)
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

	feedFollow, err := cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("handler_create_feed CreateFeedFollow error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateFeedResponse{
		Feed:        feed,
		Feed_follow: feedFollow,
	})
}
