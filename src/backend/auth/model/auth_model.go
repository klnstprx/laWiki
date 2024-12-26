package model

type User struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Email         string    `json:"email" bson:"email"`
	Name          string    `json:"name" bson:"name"`
	GivenName     string    `json:"given_name" bson:"given_name"`
	FamilyName    string    `json:"family_name" bson:"family_name"`
	Picture       string    `json:"picture" bson:"picture"`
	Locale        string    `json:"locale" bson:"locale"`
	EmailVerified bool      `json:"email_verified" bson:"email_verified"`
	Role          string    `json:"role" bson:"role"`
	Valoration    []float64 `json:"valoration" bson:"valoration"`
	Notifications []string  `json:"notifications" bson:"notifications"`
	EnableMails   bool      `json:"enable_mails" bson:"enable_mails"`
}
