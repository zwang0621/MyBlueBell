package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
)

// 处理注册请求的函数
func SignUpHandler(c *gin.Context) {

	//1.获取参数和参数校验
	var p models.ParamSignUp
	//ShouldBind只能检验数据类型对不对，请求的格式正确与否
	if err := c.ShouldBind(&p); err != nil {
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
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
