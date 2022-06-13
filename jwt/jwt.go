package jwt

import (
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"time"

	"github.com/golang-jwt/jwt"
)

var SigningKey = []byte(global.Conf.Jwt.SigningKey)

type Claims struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`

	jwt.StandardClaims
}

// GenerateToke 生成token
func GenerateToke(username string, uid int64) (string, error) {

	// 创建一个用户的声明
	claims := Claims{
		uid,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60, // 过期时间
			Issuer:    username,                  // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的key签名并获得完整的编码后的字符串token
	return token.SignedString(SigningKey)
}

// VerifyToken 验证token
func VerifyToken(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	return err
}

func GetUidByToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})

	if err != nil {
		return 0, err
	}
	// 类型断言，拿到解析后的claims

	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		uid := claims.Uid
		return uid, nil
	}
	return 0, errors.New("invalid token")

}
