package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type createUserParams struct {
	Name string `json:"name"`
}

func (cfg *apiConfig) handler_create_user(w http.ResponseWriter, r *http.Request) {
	reqParams, err := decondeJSON[createUserParams](r)

	if err != nil {
		fmt.Printf("handler_create_user decondeJSON %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if reqParams.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name is not provided")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      reqParams.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("handler_create_user db error %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
