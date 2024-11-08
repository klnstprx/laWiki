package model

import "time"

type Version struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	Editor    string    `json:"editor" bson:"editor"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	EntryID   string    `json:"entry_id" bson:"entry_id"`
}
