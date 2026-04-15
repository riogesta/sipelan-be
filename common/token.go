package common

import (
	"fmt"
	"log"
	"os"
	"sipelan/models"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretKey []byte
	once         sync.Once
)

func getJWTSecretKey() []byte {
	once.Do(func() {
		secret := os.Getenv("JWT_SECRET_KEY")
		if secret == "" {
			log.Fatal("JWT_SECRET_KEY environment variable is not set")
		}
		jwtSecretKey = []byte(secret)
	})
	return jwtSecretKey
}

func CreateToken(ID uint) (string, error) {

	claims := models.AuthClaims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 Jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(getJWTSecretKey())
	return ss, err
}

func ValidateToken(tokenString string) (*models.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecretKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.AuthClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
