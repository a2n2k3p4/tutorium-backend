package models

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type LoginRequestDoc struct {
	Username       string `json:"username" example:"b6610505511"`
	Password       string `json:"password" example:"mySecretPassword"`
	ProfilePicture string `json:"profile_picture,omitempty" example:"<base64-encoded-image>"`
	FirstName      string `json:"first_name" example:"Alice"`
	LastName       string `json:"last_name" example:"Smith"`
	Gender         string `json:"gender" example:"Female"`
	PhoneNumber    string `json:"phone_number" example:"+66912345678"`
}

type LoginResponseDoc struct {
	Token string  `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserDoc `json:"user"`
}
