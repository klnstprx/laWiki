package model

import "time"

type Entry struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Authors   []string  `json:"authors" bson:"authors"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
