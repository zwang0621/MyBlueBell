package mysql

import (
	"web_app/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertionFailed
		return
	}
	return
}

func GetCommentListByIDs(ids []string) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time
	from comment
	where comment_id in (?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids) ////用于构建带有 IN (...) 语句的 SQL 查询 —— 因为 Go 的 database/sql 包不支持直接使用 slice 作为参数传入 SQL 的 IN 子句
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindVar 的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&commentList, query, args...)
	return
}
