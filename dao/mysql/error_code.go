package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在！")
	ErrorInvalidPassword = errors.New("密码错误！")
	ErrorInvalidID       = errors.New("无效的id! ")
	ErrorInsertionFailed = errors.New("插入数据失败")
)
