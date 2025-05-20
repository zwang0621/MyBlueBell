package models

const (
	Ordertime  = "time"
	Orderscore = "score"
)

// 定义注册请求的参数结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" `
}

// 定义登录请求的参数结构体
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`              //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=-1 0 1"` //赞成票1还是反对票-1取消投票0
}

// ParamPostlist 帖子列表
type ParamPostlist struct {
	Limit       int64  `json:"limit" form:"limit"`                 //每页的数据量
	Offset      int64  `json:"offset" form:"offset"`               //起始页码
	Order       string `json:"order" form:"order" example:"score"` //排序的依据，帖子创建时间或者帖子的点赞量
	CommunityID int64  `json:"community_id" form:"community_id"`   //可以为空
}
