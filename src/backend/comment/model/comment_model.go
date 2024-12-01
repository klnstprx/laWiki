package model

import "time"

type Comment struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	Rating    int       `json:"rating" bson:"rating"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Author    string    `json:"author" bson:"author"`
	VersionID string    `json:"version_id" bson:"version_id"`
	EntryID   string    `json:"entry_id" bson:"entry_id"`
}
