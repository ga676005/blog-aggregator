package main

import (
	"time"

	"github.com/ga676005/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	UserID        uuid.UUID  `json:"userID"`
	LastFetchedAt *time.Time `json:"lastFetchedAt"`
}

func databaseFeedToFeed(f database.Feed) Feed {
	var lastFetchedAt *time.Time
	if f.LastFetchedAt.Valid {
		lastFetchedAt = &f.LastFetchedAt.Time
	}

	feed := Feed{
		ID:            f.ID,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
		Name:          f.Name,
		URL:           f.Url,
		UserID:        f.UserID,
		LastFetchedAt: lastFetchedAt,
	}

	return feed
}
