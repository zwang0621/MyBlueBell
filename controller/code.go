package controller

type ResCode int64

const (
	CodeSuccess            ResCode = 1010 + iota
	CodeInvalidParam               //1011
	CodeUserExist                  //1012
	CodeUserNotExist               //1013
	CodeInvalidPassword            //1014
	CodeServerBusy                 //1015
	CodeInvalidToken               //1016
	CodeNeedLogin                  //1017
	CodeAccessTokenExpired         //1018
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "请求参数错误",
	CodeUserExist:          "用户已存在",
	CodeUserNotExist:       "用户不存在",
	CodeInvalidPassword:    "无效密码",
	CodeServerBusy:         "服务繁忙",
	CodeInvalidToken:       "无效的token",
	CodeNeedLogin:          "需要登陆",
	CodeAccessTokenExpired: "AccessToken已过期",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
