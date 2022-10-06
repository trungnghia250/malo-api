package model

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTProfileClaims struct {
	TokenContent
	jwt.StandardClaims
}

type UserTokenClaims struct {
	Token      string  `json:"token" bson:"token"`
	UserID     string  `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreateTime float64 `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Role       string  `json:"role" bson:"role"`
}

type TokenContent struct {
	AccountType string `json:"typ,omitempty"`
	Email       string `json:"eoc,omitempty"`
	Name        string `json:"noc,omitempty"`
	Role        string `json:"role,omitempty" bson:"role,omitempty"`
}
