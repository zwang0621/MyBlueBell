package controller

import (
	"fmt"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 评论

// CommentHandler 创建评论
func CommentHandler(c *gin.Context) {
	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		fmt.Println(err)
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 生成评论ID
	commentID := snowflake.GenID()
	// 获取作者ID，当前请求的UserID
	userID, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	comment.CommentID = uint64(commentID)
	comment.AuthorID = uint64(userID)

	// 创建评论
	if err := mysql.CreateComment(&comment); err != nil {
		zap.L().Error("mysql.CreateComment(&comment) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// CommentListHandler 评论列表
func CommentListHandler(c *gin.Context) {
	ids, ok := c.GetQueryArray("ids")
	if !ok {
		ResponseError(c, CodeInvalidParam)
		return
	}
	posts, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, posts)
}
