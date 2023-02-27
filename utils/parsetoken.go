package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/wonderivan/logger"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// SECRET 加密因子
const SECRET = "yal"

// ParseToken 解析token
func ParseToken(tokenStr string) (*Claims, error) {
	//使用jwt.ParseWithClaims方法解析token，这个token是前端传给我们的,获得一个*Token类型的对象
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	//处理token解析后的各种错误
	if err != nil {
		logger.Error("parse token failed", err)
		//处理token解析后的各种报错情况
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
	}
	//转换成*CustomClaims类型并返回
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token解析失败" + err.Error())
}
