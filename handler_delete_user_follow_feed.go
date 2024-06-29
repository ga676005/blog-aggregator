package main

import (
	"fmt"
	"net/http"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

func (cfg *apiConfig) handlerDeleteUserFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	raw := r.PathValue("feedFollowID")
	if raw == "" {
		respondWithError(w, http.StatusInternalServerError, "missing id")
		return
	}

	feedFollowID, err := uuid.Parse(raw)
	if err != nil {
		fmt.Printf("handlerDeleteUserFollowFeed uuid.Pars error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.Background()
	err = cfg.DB.DeleteFeedFollow(ctx, feedFollowID)
	if err != nil {
		fmt.Printf("handlerDeleteUserFollowFeed DeleteFeedFollow error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}
