package models

import "time"

/*
字段 ParentID 表示 该评论的父评论ID,也就是说：
如果 ParentID == 0,那么这个评论是对帖子的一级评论；
如果 ParentID != 0,那么这个评论是对某条其他评论的回复，其值等于被回复评论的 CommentID

假设有一个帖子 PostID = 1:
用户A评论了帖子,CommentID = 1001, ParentID = 0
用户B回复了用户A的评论,CommentID = 1002, ParentID = 1001
用户C又回复了用户B的评论,CommentID = 1003, ParentID = 1002

这样的结构就形成了一个“评论树”或“楼中楼”的层级关系
*/
type Comment struct {
	PostID     uint64    `db:"post_id" json:"post_id"`
	ParentID   uint64    `db:"parent_id" json:"parent_id"`
	CommentID  uint64    `db:"comment_id" json:"comment_id"`
	AuthorID   uint64    `db:"author_id" json:"author_id"`
	Content    string    `db:"content" json:"content"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
}
