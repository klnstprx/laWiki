package model

type Version struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Content  string `json:"content" bson:"content"`
	Author   string `json:"author" bson:"author"`
	NVersion string `json:"nVersion" bson:"nVersion"`
}
