package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(userId int64, tokenTTL time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userId
	claims["exp"] = time.Now().Add(tokenTTL).Unix()
	tokenHash, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenHash, nil
}
