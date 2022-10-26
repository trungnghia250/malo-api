package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/service/model"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
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

type AuthenticationConfig struct {
	Key                   string
	AuthorizationDatabase *mongo.Database
	Claims                interface{}
	Language              string
	JWTExp                int
	ValidateExpType       string
	JWTConfig             middleware.JWTConfig
}

func VerifyJWTToken(config AuthenticationConfig, rawToken string) (*jwt.Token, *model.JWTProfileClaims, error) {
	t := reflect.ValueOf(config.JWTConfig.Claims).Type().Elem()
	claims := reflect.New(t).Interface().(jwt.Claims)
	token, err := jwt.ParseWithClaims(rawToken, claims, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.JWTConfig.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.JWTConfig.SigningKey, nil
	})
	if err != nil && err.(*jwt.ValidationError).Errors == jwt.ValidationErrorMalformed {
		return nil, nil, errors.New("Not Authorization")
	}
	claim := token.Claims.(*model.JWTProfileClaims)
	if !token.Valid {
		// TODO Allow expired token for now
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired {
				// allow expired token
			} else {
				return nil, claim, errors.New("Not Authorized")
			}
		}
	}
	return token, claim, nil
}
