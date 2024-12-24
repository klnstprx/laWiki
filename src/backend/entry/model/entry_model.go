package model

import (
	"time"
)

type Entry struct {
	ID               string                       `json:"id" bson:"_id,omitempty"`
	Title            string                       `json:"title" bson:"title"`
	Author           string                       `json:"author" bson:"author"`
	CreatedAt        time.Time                    `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time                    `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	WikiID           string                       `bson:"wiki_id" json:"wiki_id"`
	TranslatedFields map[string]map[string]string `json:"translatedFields,omitempty" bson:"translatedFields,omitempty"`
	SourceLang       string                       `json:"sourceLang,omitempty" bson:"sourceLang,omitempty"`
}
