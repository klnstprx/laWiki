package model

type Version struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Content  string `json:"content" bson:"content"`
	Editor   string `json:"editor" bson:"editor"`
	NVersion string `json:"n_version" bson:"n_version"`
	EntryID  string `json:"entry_id" bson:"entry_id"`
}
