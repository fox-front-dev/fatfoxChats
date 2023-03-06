package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Cla struct {
	Id uint64
	jwt.RegisteredClaims
}

var ChatsSigned = []byte("fatfox")

// 获取token
func GetToken(userId uint64) (string, error) {
	cla := Cla{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// fmt.Println(cla)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	return token.SignedString(ChatsSigned)
}

// 解析token
func ParseToken(token string) (*Cla, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Cla{}, func(t *jwt.Token) (interface{}, error) {
		return ChatsSigned, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := jwtToken.Claims.(*Cla); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
