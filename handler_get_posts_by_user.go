package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "20"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.Background()
	feedsFollow, err := cfg.DB.GetUserFeedsFollow(ctx, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedsId := make([]uuid.UUID, 0)
	for _, v := range feedsFollow {
		feedsId = append(feedsId, v.FeedID)
	}

	if len(feedsId) == 0 {
		respondWithJSON(w, http.StatusOK, []any{})
	}

	posts, err := cfg.DB.GetPostsByUser(ctx, database.GetPostsByUserParams{FeedIds: feedsId, Limit: int32(limit)})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}
