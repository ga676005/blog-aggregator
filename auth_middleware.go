package main

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ga676005/blog-aggregator/internal/database"
	"golang.org/x/net/context"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		reqApiKey := strings.TrimPrefix(authHeader, "ApiKey ")
		if reqApiKey == "" {
			respondWithError(w, http.StatusBadRequest, "no apiKey")
			return
		}

		ctx := context.Background()
		user, err := cfg.DB.GetUserByKey(ctx, reqApiKey)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "user not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		handler(w, r, user)
	})
}
