package model

type Media struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	PublicID  string `json:"publicId" bson:"publicId"`
	UploadUrl string `json:"uploadUrl" bson:"uploadUrl"`
}
