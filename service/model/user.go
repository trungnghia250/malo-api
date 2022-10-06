package model

type User struct {
	UserID       string `bson:"user_id" json:"user_id"`
	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	Password     string `bson:"password" json:"password"`
	Role         string `bson:"role" json:"role"`
	Token        string `json:"token" bson:"token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
