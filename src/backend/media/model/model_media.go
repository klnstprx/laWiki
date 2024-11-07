package model

type Media struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	UploadUrl string `json:"uploadUrl" bson:"uploadUrl"`
	AssetUrl  string `json:"assetUrl" bson:"assetUrl"`
}
