package model

import "time"

type Wiki struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Category    string    `json:"category" bson:"category"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	MediaID     string    `json:"media_id,omitempty" bson:"media_id,omitempty"`
}
