package controller

import (
	"strconv"
	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 跟社区相关的
func CommunityHandler(c *gin.Context) {
	//查询所有的社区（communityid，communityname）以列表的形式返回
	//获取所有社区信息不需要进行参数校验，直接从 Logic 层拉取信息即可
	//并不是所有的get方法都不需要进行参数校验，比如说通过query传参
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	communityIDstr := c.Param("id") //c.Param获取路径参数
	id, err := strconv.ParseInt(communityIDstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id获取社区的详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
