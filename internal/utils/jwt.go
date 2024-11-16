package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JWTSECRET = "secret"

func GenerateToken(userId int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(JWTSECRET))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func jwtParser(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSECRET), nil
	})
}
