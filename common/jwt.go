package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	DefaultExpire int = 2 * 3600
	jwtSecret         = []byte("tip$%^$%*&tok")
)

type (
	UserTokenClaims struct {
		Username string                 `json:"username"`
		Password string                 `json:"password"`
		AddData  map[string]interface{} `json:"addData"`
		jwt.StandardClaims
	}
	JwtOptions struct {
		Expire        int //second
		JwtSecret     []byte
		UseJSONNumber bool
		AddData       map[string]interface{}
	}
	JwtOption func(options *JwtOptions)
)

//解析 UserTokenClaims
func ParseJWTToken(token string, options ...JwtOption) (*UserTokenClaims, error) {
	option := NewJwtOptions(options...)
	parser := jwt.Parser{UseJSONNumber: option.UseJSONNumber}
	tokenClaims, err := parser.ParseWithClaims(token, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return option.JwtSecret, nil
	})
	if tokenClaims != nil {
		if claim, ok := tokenClaims.Claims.(*UserTokenClaims); ok && tokenClaims.Valid {
			return claim, nil
		}
	}

	return nil, err
}

func GenerateToken(username, password string, options ...JwtOption) (string, error) {
	now := time.Now()
	option := NewJwtOptions(options...)

	expireTime := now.Add(time.Second * time.Duration(option.Expire))
	claims := UserTokenClaims{
		Username: username,
		Password: password,
		AddData:  option.AddData,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "jwt",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(option.JwtSecret)
	return token, err
}

func WithExpire(expire int) func(options *JwtOptions) {
	return func(options *JwtOptions) {
		options.Expire = expire
	}
}

func WithJWTSecret(secret []byte) func(options *JwtOptions) {
	return func(options *JwtOptions) {
		options.JwtSecret = secret
	}
}

func WithAddData(data map[string]interface{}) func(options *JwtOptions) {
	return func(options *JwtOptions) {
		options.AddData = data
	}
}

func WithUseJSONNumber(useJsonNumber bool) func(options *JwtOptions) {
	return func(options *JwtOptions) {
		options.UseJSONNumber = useJsonNumber
	}
}

func NewJwtOptions(options ...JwtOption) *JwtOptions {
	option := &JwtOptions{
		Expire:        DefaultExpire,
		AddData:       make(map[string]interface{}),
		JwtSecret:     jwtSecret,
		UseJSONNumber: true,
	}
	for i := range options {
		options[i](option)
	}
	return option
}
