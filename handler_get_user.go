package main

import (
	"net/http"

	"github.com/ga676005/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handler_get_user(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}
