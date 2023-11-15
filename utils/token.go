package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmcerez0/fiber-demo/models"
)

func GetToken(user models.User) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.FullName(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
