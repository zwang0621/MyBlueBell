package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("火车叨位去")
var ErrTokenExpired = errors.New("token has expierd") //为过期的atoken定义一个特定的错误
var ErrTokenValid = errors.New("token is invalid")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// GenToken 生成access token 和 refresh token
func GenToken(userid int64, username string) (atoken, rtoken string, err error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		userid,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), //过期时间
			Issuer:    "bluebell",                                                                        // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象，加密并获得完整的编码后的字符串token
	atoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySecret)
	if err != nil {
		return "", "", err
	}

	// refresh token 不需要存任何自定义数据
	rtoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),
		Issuer:    "bluebell",
	}).SignedString(mySecret)
	if err != nil {
		return "", "", err
	}

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

// GenToken 生成JasonWebToken
func GenToken2(userid int64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		userid,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), //过期时间
			Issuer:    "bluebell",                                                                        // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	// 直接使用标准的Claim则可以直接使用Parse方法
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error)封装成了一个函数keyFunc
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		//检查错误是否为jwt验证错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			//检查验证错误是否明确是由token过期导致的
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		//其他错误返回一个通用的无效token
		return nil, ErrTokenValid
	}

	if token != nil && token.Valid {
		return mc, nil
	}
	// // 对token对象中的Claim进行类型断言
	// if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
	// 	return claims, nil
	// }
	return nil, ErrTokenValid
}

// RefreshToken 刷新 accesstoken
func RefreshToken(rtoken string) (newAtoken, newRtoken string, err error) {
	var claims MyClaims
	// refresh token无效直接返回
	if _, err = jwt.ParseWithClaims(rtoken, &claims, keyFunc); err != nil {
		return "", "", err
	}
	//生成新的atoken 和 rtoken
	//这里生成新的rtoken 安全性更高，防止rtoken泄露
	//如果只返回atoken，虽然实现更简单，但是安全性会降低
	return GenToken(claims.UserID, claims.Username)
}
