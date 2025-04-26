package models

import "time"

type NewsItems struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
}
