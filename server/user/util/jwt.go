package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secret_key = "miolyn"

type JWT struct {
	SigningKey []byte //jwt签名
}

type JWTClaims struct {
	//继承jwt.MapClaims 的方法
	jwt.MapClaims
	UserID    string //用户名
	ExpiredAt int64  //过期时间
	CreatedAt int64  //生效时间
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(secret_key),
	}
}

//创建Token
func (j *JWT) GenerateToken(claims JWTClaims) (string, error) {
	//设置过期时间
	claims.ExpiredAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims.CreatedAt = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解析Token
func (j *JWT) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
	}
	return nil, errors.New("Token is not valid")

}
