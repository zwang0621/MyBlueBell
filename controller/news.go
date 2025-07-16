package controller

import (
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewsTrendingHandler 获取热点新闻
func NewsTrendingHandler(c *gin.Context) {
	p := &models.ParamNewsTrending{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("NewsTrendingHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetBaiduNewsTrending(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
