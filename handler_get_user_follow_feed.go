package main

import (
	"fmt"
	"net/http"

	"github.com/ga676005/blog-aggregator/internal/database"
	"golang.org/x/net/context"
)

func (cfg *apiConfig) handlerGetUserFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	feedsFollow, err := cfg.DB.GetUserFeedsFollow(ctx, user.ID)
	if err != nil {
		fmt.Printf("handlerGetUserFollowFeed GetUserFeedsFollow error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果 feedsFollow 沒東西會是 nil slice，JSON marshall 後會變成 nil
	// 所以要回傳 [] 的話要特別處理
	if feedsFollow == nil {
		respondWithJSON(w, http.StatusOK, []any{})
		return
	}

	respondWithJSON(w, http.StatusOK, feedsFollow)
}
