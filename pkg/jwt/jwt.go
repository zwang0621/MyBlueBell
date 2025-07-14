package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("火车叨位去")

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
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 accesstoken
func RefreshToken(atoken, rtoken string) (newAtoken, newRtoken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rtoken, keyFunc); err != nil {
		return "", "", err
	}

	//从旧access token中解析出claims数据 解析出payload 负载信息
	var claims MyClaims
	_, err = jwt.ParseWithClaims(atoken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	//当accesstoken是过期错误，并且refreshtoken没有过期就创建一个新的accesstoken
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
