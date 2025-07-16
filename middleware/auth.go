package middleware

import (
	"errors"
	"strings"
	"web_app/controller"
	"web_app/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			// c.JSON(http.StatusOK, gin.H{
			// 	"code": 2003,
			// 	"msg":  "请求头中auth为空",
			// })
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			// c.JSON(http.StatusOK, gin.H{
			// 	"code": 2004,
			// 	"msg":  "请求头中auth格式有误",
			// })
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// 如果 access token 过期，我们返回一个特定的、约定好的状态码。
				// 客户端将拦截这个状态码并触发Token刷新流程。
				controller.ResponseErrorWithMsg(c, controller.CodeAccessTokenExpired, "access token已过期")
				c.Abort()
				return
			}

			// 对于所有其他Token错误（例如，签名无效、格式错误），
			// 返回通用的“无效Token”错误。
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, "无效的token")
			c.Abort()
			return
		}
		// 将当前请求的userid信息保存到请求的上下文c上
		c.Set(controller.ContextUserIdKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIdKey)来获取当前请求的用户信息
	}
}
