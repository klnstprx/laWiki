package model

import "time"

type Version struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	Editor    string    `json:"editor" bson:"editor"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	EntryID   string    `json:"entry_id" bson:"entry_id"`
	Address   string    `json:"address" bson:"address"`
	MediaIDs  []string  `json:"media_ids,omitempty" bson:"media_ids,omitempty"`
}
