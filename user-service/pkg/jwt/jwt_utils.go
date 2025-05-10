package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"time"
)

type Utils interface {
	CreateToken(userId int, role string) (string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

const jwtTokenExpTime = 60 * time.Minute

type jwtUtils struct {
	issuer    string
	secretKey string
}

func (u *jwtUtils) CreateToken(userId int, role string) (string, error) {
	expireTime := time.Now().Add(jwtTokenExpTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expireTime,
		"role":   role,
		"iss":    u.issuer,
	})
	tokenString, err := token.SignedString([]byte(u.secretKey))
	if err != nil {
		log.Println("Jwt service, create access token err :", err)
		return "", errors.New("internal server error")
	}
	return tokenString, nil
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
	return &jwtUtils{
		issuer:    os.Getenv("ISSUER"),
		secretKey: os.Getenv("SECRET_KEY"),
	}
}
