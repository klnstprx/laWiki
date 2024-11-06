package model

import "time"

type Comment struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	Rating    int       `json:"rating" bson:"rating"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
