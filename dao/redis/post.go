package redis

import (
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis"
)

func getIDsFormKey(key string, Offset, Limit int64) ([]string, error) {
	//确定查询的索引起始点
	start := (Offset - 1) * Limit
	end := start + Limit - 1
	//ZRevRange按分数从大到小查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()

}
func GetPostIDsInOrder(p *models.ParamPostlist) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.Orderscore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	return getIDsFormKey(key, p.Offset, p.Limit)
}

// 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	// 	//查找key中分数是1的元素的数量，即每篇帖子的赞成票的数量
	// 	v := rdb.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	// keys := make([]string, 0, len(ids))

	//使用pipeline一次发送多条命令，减少rtt
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	data = make([]int64, 0, len(cmders))
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostlist) ([]string, error) {
	//使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据
	orderkey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.Orderscore {
		orderkey = getRedisKey(KeyPostScoreZSet)
	}
	ckey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID))) //社区的key
	//利用缓存key减少zinterstore执行的次数
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := rdb.Pipeline()
		//用 ZINTERSTORE 把社区下的帖子集合（Set）全局排序的 ZSet（时间/得分）做一个交集，生成一个新的 ZSet：
		//key 是上面新拼的 key，分数来自原排序 ZSet，用 MAX 保留最大的分数（通常其实只有一个来源）
		//得到一个："某社区下的帖子，按时间/得分排序的新 ZSet"
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderkey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话就直接根据key查询ids
	return getIDsFormKey(key, p.Offset, p.Limit)
}
