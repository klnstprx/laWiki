package model

import (
	"time"
)

type Entry struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Author    string    `json:"author" bson:"author"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	WikiID    string    `bson:"wiki_id" json:"wiki_id"`
}
