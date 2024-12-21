package model

type User struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	GivenName     string  `json:"given_name"`
	FamilyName    string  `json:"family_name"`
	Picture       string  `json:"picture"`
	Locale        string  `json:"locale"`
	VerifiedEmail bool    `json:"verified_email"`
	Role          string  `json:"role"`
	Valoration    float64 `json:"valoration"`
}
