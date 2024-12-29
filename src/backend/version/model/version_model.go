package model

import "time"

// Translation representa una traducción de una entrada
type Translation struct {
	Language string `json:"language" bson:"language"`
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
}

type Version struct {
	ID               string                       `json:"id" bson:"_id,omitempty"`
	Content          string                       `json:"content" bson:"content"`
	TranslatedFields map[string]map[string]string `json:"translatedFields,omitempty" bson:"translatedFields,omitempty"`
	SourceLang       string                       `json:"sourceLang" bson:"sourceLang"`
	Editor           string                       `json:"editor" bson:"editor"`
	CreatedAt        time.Time                    `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time                    `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	EntryID          string                       `json:"entry_id" bson:"entry_id"`
	Address          string                       `json:"address" bson:"address"`
	MediaIDs         []string                     `json:"media_ids,omitempty" bson:"media_ids,omitempty"`
}
