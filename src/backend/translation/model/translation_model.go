package model

import "time"

// Translation representa una traducción de una entrada o versión
type Translation struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	SourceID       string    `json:"sourceId" bson:"source_id"`
	SourceLang     string    `json:"sourceLang" bson:"source_lang"`
	TargetLang     string    `json:"targetLang" bson:"target_lang"`
	SourceText     string    `json:"sourceText" bson:"source_text"`
	TranslatedText string    `json:"translatedText" bson:"translated_text"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}
