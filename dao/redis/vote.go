package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeat     = errors.New("不许重复投票")
)

func CreatePost(postID, community_id int64) error {
	pipeline := rdb.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.FormatInt(postID, 10),
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		//初始的分数仍然与时间相关联，这与直觉相符，越新的帖子分数应该越高，使得新的帖子尽可能靠前显示
		Score:  float64(time.Now().Unix()),
		Member: strconv.FormatInt(postID, 10),
	})

	//补充 把帖子id加到社区的set
	pipeline.SAdd(getRedisKey(KeyCommunitySetPrefix+strconv.Itoa(int(community_id))), postID)
	_, err := pipeline.Exec()

	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取发布时间
	post_time := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-post_time > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}
	//2和3需要放到一个pipeline事务里面

	//2.更新帖子的分数
	//先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	//不能重复投相同的票
	if value == ov {
		return ErrorVoteRepeat
	}
	var symbol float64
	if value > ov {
		symbol = 1
	} else {
		symbol = -1
	}
	diff := math.Abs(ov - value)
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), symbol*diff*scorePerVote, postID)
	//3.记录用户为该则帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		//如果用户之前已经使用了 ZAdd 添加相同的 Member 到 Sorted Set，那么本次 ZAdd 将会把之前的 Member 和 Score 覆盖掉
		//也就是说如果用户之前对该帖子投了up，但是这次投了down，那么redis KeyPostVotedZSetPrefix就只会记录最新的
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
