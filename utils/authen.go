package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/service/model"
	"time"
)

func GenToken(tokenContent model.TokenContent) (string, error) {
	claims := &model.JWTProfileClaims{
		TokenContent:   tokenContent,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.Config.Encryption.JWTExp)).Unix(), Issuer: "malo", IssuedAt: time.Now().Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(fmt.Sprintf("%s", config.Config.Encryption.JWTSecret)))
}

func GenRefreshToken() string {
	return GetMD5Hash(uuid.NewV4().String())
}

func GetMD5Hash(text string) string {
	harsher := md5.New()
	harsher.Write([]byte(text))
	return hex.EncodeToString(harsher.Sum(nil))
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(config.Config.Encryption.JWTSecret), nil
	})
}
