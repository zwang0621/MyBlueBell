package redis

//redis key
//redis key尽量用命名空间的方式区分不同的key

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //ZSet 帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //ZSet 帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" //ZSet 记录用户及投票的类型,参数是post id
	KeyCommunitySetPrefix  = "community:"  //set  保存每个分区下帖子的id
)

//给rediskey加前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
