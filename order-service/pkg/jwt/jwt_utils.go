package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
)

type Utils interface {
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type jwtUtils struct {
}

func (*jwtUtils) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("Parse token err : %v", err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Parse token err : token is invalid")
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func NewJwtUtils() Utils {
	return &jwtUtils{}
}
