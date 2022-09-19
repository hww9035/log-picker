package other

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// 过期时间
const tokenExpireDuration = time.Hour * 2

var (
	mySecret = []byte("hww9035")
	issuer   = "hww9035"
)

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	c := MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    issuer,
		},
		username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
