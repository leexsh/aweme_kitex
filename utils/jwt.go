package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JwtKey            = "my_aweme_kitex"
	TokenExpire int64 = 3600 * 12
)

type UserClaim struct {
	jwt.StandardClaims
	Id   int64
	Name string
}

func GenerateToken(id int64, name string, second int64) (string, error) {
	uc := UserClaim{
		Id:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyzeToke(token string) (*UserClaim, error) {
	uc := &UserClaim{}
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, errors.New("token is invalid")
	}
	return uc, err
}
