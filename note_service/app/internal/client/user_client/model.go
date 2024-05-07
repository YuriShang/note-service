package user_client

import "github.com/google/uuid"

type User struct {
	UUID            uuid.UUID `json:"id" bson:"_id"`
	Username        string    `json:"username" bson:"username"`
	RegisterTime    string    `json:"register_time" bson:"register_time"`
	PasswordSetTime string    `json:"password_set_time" bson:"password_set_time"`
}

type Token struct {
	AccessToken string `json:"access_token" bson:"access_token"`
	TokenType   string `json:"token_type" bson:"token_type"`
}
