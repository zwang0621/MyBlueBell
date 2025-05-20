package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code":10001,      //程序中的错误码
	"msg": xx,         //提示信息
	"data":{},         //数据
}
*/

func ResponseError(ctx *gin.Context, c ResCode) {

	ctx.JSON(http.StatusOK, gin.H{
		"code": c,
		"msg":  c.Msg(),
		"data": nil,
	})

}

func ResponseErrorWithMsg(ctx *gin.Context, c ResCode, msg interface{}) {

	ctx.JSON(http.StatusOK, gin.H{
		"code": c,
		"msg":  msg,
		"data": nil,
	})

}

func ResponseSuccess(ctx *gin.Context, data interface{}) {

	ctx.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
		"data": data,
	})

}
