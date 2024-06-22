package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type requestParam struct {
	Name string `json:"name"`
}

func (cfg *apiConfig) handler_create_user(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)
	params := requestParam{}
	err := decorder.Decode(&params)
	if err != nil {
		fmt.Printf("handler_create_user couldn't decode params %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name is not provided")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
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
