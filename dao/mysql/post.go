package mysql

import (
	"strings"
	"web_app/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// CreatePost创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID根据id查询单个帖子详情
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	if err != nil {
		zap.L().Error("Failed to execute query with param", zap.Error(err))
	}
	return
}

// GetPostList查询帖子列表
func GetPostList(offset, limit int64) (posts []*models.Post, err error) {
	/*
		db.Select 可以正确地向一个 nil 切片追加元素，所以从功能上讲，make 并不是必须的
		使用 make 创建具有初始容量的切片，通常是为了性能优化，通过预分配内存来减少后续追加操作中的扩容开销
		同时，它也可能是一种代码风格，让切片的创建和初始状态（长度为 0）更加明确
	*/
	posts = make([]*models.Post, 0, 2) //长度，容量
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post ORDER BY create_time DESC limit ?,?`
	db.Select(&posts, sqlStr, (offset-1)*limit, limit)
	return
}

// GetPostListByIds根据给定的id列表查询帖子数据
func GetPostListByIds(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id in (?) order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

// GetPostListTotalCount 根据关键词查询帖子列表总数
func GetPostListTotalCount(p *models.ParamPostList) (count int64, err error) {
	// 根据帖子标题或者帖子内容模糊查询帖子列表总数
	sqlStr := `SELECT count(post_id)
	FROM post
	WHERE (title LIKE ? OR content LIKE ?)
	`

	// %keyword%
	p.Search = "%" + p.Search + "%"
	err = db.Get(&count, sqlStr, p.Search, p.Search)
	return
}

// GetPostListByKeywords 根据关键词查询帖子列表
func GetPostListByKeywords(p *models.ParamPostList) (posts []*models.Post, err error) {
	// 根据帖子标题或者帖子内容模糊查询帖子列表
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time
	FROM post
	WHERE (title LIKE ? OR content LIKE ?)
	ORDER BY create_time DESC
	LIMIT ?,?
	`
	// %keyword%
	p.Search = "%" + p.Search + "%"
	posts = make([]*models.Post, 0, 10)
	page := p.Offset
	if page < 1 {
		page = 1 // 保证页码至少为1
	}
	offset := (page - 1) * p.Limit // 标准的分页计算公式
	err = db.Select(&posts, sqlStr, p.Search, p.Search, offset, p.Limit)
	return
}
