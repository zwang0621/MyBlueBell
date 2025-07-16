package logic

import (
	"strconv"
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

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
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

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
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
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
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

// PostSearch 搜索业务-搜索帖子
func PostSearch(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes
	// 根据搜索条件去mysql查询符合条件的帖子列表总数
	total, err := mysql.GetPostListTotalCount(p)
	if err != nil {
		return nil, err
	}
	res.Page.Total = total
	// 1、根据搜索条件去mysql分页查询符合条件的帖子列表
	posts, err := mysql.GetPostListByKeywords(p)
	if err != nil {
		return nil, err
	}
	// 查询出来的帖子总数可能为0
	if len(posts) == 0 {
		return &models.ApiPostDetailRes{}, nil
	}
	// 2、查询出来的帖子id列表传入到redis接口获取帖子的投票数
	ids := make([]string, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, strconv.Itoa(int(post.ID)))
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Size = p.Limit
	res.Page.Page = p.Offset
	// 3、拼接数据
	res.List = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", uint64(post.AuthorID)),
				zap.Error(err))
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				//这会在日志中生成一条格式化的信息，如：
				//ERROR  mysql.GetCommunityByID() failed  {"community_id": 12345, "error": "sql: no rows in result set"}
				//为了在日志中清晰地标记出是哪一个社区的查询失败了，方便排查问题和定位具体数据。
				zap.Uint64("community_id", uint64(post.CommunityID)),
				zap.Error(err))
		}

		authorName := "未知用户" // 提供一个默认值
		if user != nil {
			authorName = user.Username
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      authorName,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}
