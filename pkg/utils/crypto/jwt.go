package crypto

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(login, secretKey string) (string, error) {
	signingKey := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString(signingKey)
}
