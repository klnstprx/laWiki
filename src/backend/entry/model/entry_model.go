package model

import (
	"time"
)

type Entry struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Authors   []string  `json:"authors" bson:"authors"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
