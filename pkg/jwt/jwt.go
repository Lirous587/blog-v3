package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenParams struct {
	Payload  interface{}
	Secret   []byte
	Duration time.Duration
}

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type MyClaims[T any] struct {
	PayLoad T `json:"payLoad"`
	jwt.RegisteredClaims
}

func GenToken[T any](data *JWTTokenParams) (string, error) {
	payload, ok := data.Payload.(T)
	if !ok {
		return "", errors.New("invalid payload type")
	}

	claims := &MyClaims[T]{
		PayLoad: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(data.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Lirous-blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(data.Secret)
}

func ParseToken[T any](tokenString string, secret []byte) (myClaims *MyClaims[T], needRefresh bool, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			if claims, ok := token.Claims.(*MyClaims[T]); ok {
				return claims, true, nil
			}
		}
		return nil, false, err
	}

	if claims, ok := token.Claims.(*MyClaims[T]); ok {
		return claims, false, nil
	}

	return nil, false, errors.New("invalid token")
}
