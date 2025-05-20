package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

const ContextUserIdKey = "userID"

// GetCurrentUser获取当前登录的用户id
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIdKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPageInfo(c *gin.Context) (int64, int64) {
	offset_str := c.Query("offset")
	limit_str := c.Query("limit")
	var (
		offset int64
		limit  int64
		err    error
	)
	offset, err = strconv.ParseInt(offset_str, 10, 64)
	if err != nil {
		offset = 0
	}
	limit, err = strconv.ParseInt(limit_str, 10, 64)
	if err != nil {
		limit = 10
	}
	return offset, limit
}
