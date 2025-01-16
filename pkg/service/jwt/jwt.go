package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type InfoData[T any] struct {
	Core T
	jwt.RegisteredClaims
}

func Create[T any](info T, duration time.Duration, SECRET string, Issuer string) (string, error) {
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		InfoData[T]{
			Core: info,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
				Issuer:    Issuer,
			},
		})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func Parse[T any](token string, SECRET string) (*T, error) {
	var res InfoData[T]
	// Parse the token
	info, err := jwt.ParseWithClaims(token, &res, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := info.Claims.(*InfoData[T]); ok && info.Valid {
		return &claims.Core, nil
	} else {
		return nil, fmt.Errorf("token不合法,或者已过期")
	}
}
