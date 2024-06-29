package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type postUserFollowFeedParams struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (cfg *apiConfig) handlerPostUserFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	reqParams, err := parseRequestJSON[postUserFollowFeedParams](r)
	if err != nil {
		fmt.Printf("handlerPostUserFollowFeed parseRequestJSON error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.Background()

	_, err = cfg.DB.GetFeed(ctx, reqParams.FeedID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusInternalServerError, "no feed")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return

	}

	feedFollow, err := cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    reqParams.FeedID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("handlerPostUserFollowFeed CreateFeedFollow error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}
