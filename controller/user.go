package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
	"web_app/pkg/jwt"
)

// 处理注册请求的函数
func SignUpHandler(c *gin.Context) {

	//1.获取参数和参数校验
	var p models.ParamSignUp
	//ShouldBind只能检验数据类型对不对，请求的格式正确与否
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err)) //记录进日志
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			ResponseError(c, CodeInvalidParam) //直接替换
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": RemoveTopStruct(errs.Translate(trans)),
		// })
		ResponseErrorWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans))) //直接替换
		return
	}

	//手动对请求参数进行详细的业务规则校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	// 	//请求参数有误，直接返回响应
	// 	zap.L().Error("SignUp with invalid param") //记录进日志
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }
	fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(&p); err != nil {
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "注册失败！",
		// })
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "success",
	// })
	ResponseSuccess(c, nil)
}

// 处理登录请求的函数
func LoginHandler(c *gin.Context) {

	//1.获取参数和参数校验
	var p models.ParamLogin
	//ShouldBind只能检验数据类型对不对，请求的格式正确与否
	if err := c.ShouldBind(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err)) //记录进日志
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			ResponseError(c, CodeInvalidParam) //直接替换
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": RemoveTopStruct(errs.Translate(trans)),
		// })
		ResponseErrorWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans)))
		return
	}

	fmt.Println(p)

	//2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.login failed", zap.String("username", p.Username), zap.Error(err))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "登录失败！",
		// })
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登陆成功！",
	// })
	ResponseSuccess(c, gin.H{
		"userID":   fmt.Sprintf("%d", user.UserID), //前端js识别的最大值：id值大于1<<53-1  int64: i<<63-1
		"username": user.Username,
		"password": user.Password,
		"atoken":   user.AccessToken,
		"rtoken":   user.RefreshToken,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// 处理rtoken刷新atoken的函数
func RefreshTokenHandler(c *gin.Context) {
	// 1. 从请求体中获取 refresh_token
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果请求体解析失败，说明客户端没有按规定传递参数
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		c.Abort()
		return
	}

	// 2. 调用业务逻辑来刷新 token
	// 注意：刷新 access_token 只需要 refresh_token
	aToken, rToken, err := jwt.RefreshToken(req.RefreshToken) // 假设 RefreshToken 函数只需要 refresh_token
	if err != nil {
		// 如果刷新失败（例如 refresh_token 无效或过期）
		zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidToken, "refresh token is invalid or expired") // 使用你自己的错误响应函数
		c.Abort()
		return
	}

	// 3. 成功刷新，返回新的 token
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken, // 返回新的 refresh_token (如果实现了 Rotation)
	})
}
