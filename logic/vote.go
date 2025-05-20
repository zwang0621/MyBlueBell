package logic

import (
	"strconv"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

//投票功能

/*投票的几种情况
direction=1时，有两种情况：
	1.之前没投过票，现在投赞成
	2.之前投反对票，现在投赞成

direction=0时，有两种情况：
	1.之前投赞成票，现在取消投票
	2.之前投反对票，现在取消投票

direction=-1时，有两种情况：
	1.之前没投过票，现在投反对
	2.之前投赞成票，现在投反对

投票的限制：
每个帖子自发表之日起，一个星期之内允许用户投票，超过一个星期不允许在投票
	1.到期之后将redis中保存的赞成票数和反对票数存储到mysql中
	2.到期之后删除记录投票的redis key KeyPostVotedZSetPrefix
*/

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", p.PostID), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
