package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) getAllFeeds(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	feeds, err := cfg.DB.GetAllFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
