package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
	//3.返回
}

// 根据帖子id查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {

	//查询并拼接组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Error(err))
		return
	}
	//根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList获取帖子列表
func GetPostList(offset, limit int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(offset, limit)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamPostlist) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 3.根据id去mysql数据库查详细信息
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数，不需要在下面的循环中去查
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostlist) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	// 3.根据id去mysql数据库查详细信息
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数，不需要在下面的循环中去查
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListNew 将两个查询帖子列表接口合二为一
func GetPostListNew(p *models.ParamPostlist) (data []*models.ApiPostDetail, err error) {
	/*
		如果 CommunityID 为 0，那么默认从所有 Community 按分数排序对帖子进行查询，否则根据指定的 CommunityID 对帖子进行查询
	*/
	//查所有
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		//根据社区id查
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew() failed", zap.Error(err))
	}
	return data, err
}
