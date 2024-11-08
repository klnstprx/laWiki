package model

type Wiki struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Category    string `json:"category" bson:"category"`
}
