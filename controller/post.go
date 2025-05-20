package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（帖子的id，从url中获取）
	pid_str := c.Param("id")
	pid, err := strconv.ParseInt(pid_str, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子的数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	offset, limit := GetPageInfo(c) //这里的 getPageInfo 函数是用于获取当前页的配置信息的函数。需要配置 PageInfo 的原因是，前端页面一次可能显示不下所有帖子的信息，因此需要分页
	//获取数据
	data, err := logic.GetPostList(offset, limit)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// 升级版帖子列表接口2，根据前端传来的参数（按分数或者按创建时间）动态的去获取帖子列表
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostlist false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// 1.获取参数
	// 2.去redis查id列表
	// 3.根据id去数据库查
	p := models.ParamPostlist{
		Offset: 1,
		Limit:  10,
		Order:  models.Ordertime,
	}
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(&p) //更新：合二为一

	if err != nil {
		zap.L().Error("logic.GetPostListNew() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}
