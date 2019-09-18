package common

import (
	"time"
	"github.com/tiptok/gocomm/pkg/log"

	"github.com/dgrijalva/jwt-go"
)


type UserTokenClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var jwtSecret = []byte("123456")

//解析 UserTokenClaims
func ParseJWTToken(token string) (*UserTokenClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claim, ok := tokenClaims.Claims.(*UserTokenClaims); ok && tokenClaims.Valid {
			log.Info("ParseJWTToken:%s -> %v", token, claim)
			return claim, nil
		}
	}

	return nil, err
}

func GenerateToken(username, password string) (string, error) {
	now := time.Now()
	expireTime := now.Add(3 * time.Hour)

	claims := UserTokenClaims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "jwt",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}