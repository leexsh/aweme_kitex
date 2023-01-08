package jwt

import (
	"aweme_kitex/pkg/errno"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JwtKey            = "my_aweme_kitex"
	TokenExpire int64 = 3600 * 24 * 365 * 10

	ErrTokenExpired     = errno.WithCode(errno.TokenExpiredErrCode, "Token expired")
	ErrTokenNotValidYet = errno.WithCode(errno.TokenValidationErrCode, "Token is not active yet")
	ErrTokenMalformed   = errno.WithCode(errno.TokenInvalidErrCode, "That's not even a token")
	ErrTokenInvalid     = errno.WithCode(errno.TokenInvalidErrCode, "Couldn't handle this token")
)

type UserClaim struct {
	jwt.StandardClaims
	Id   string
	Name string
}

func GenerateToken(id, name string) (string, error) {
	uc := UserClaim{
		Id:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(TokenExpire)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyzeToken(token string) (*UserClaim, error) {
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
