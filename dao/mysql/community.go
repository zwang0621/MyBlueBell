package mysql

import (
	"database/sql"
	"web_app/models"

	"go.uber.org/zap"
)

// GetCommunityList 展示所有的社区id和社区名字
func GetCommunityList() (communityList []*models.Community, err error) {
	//为什么这里不需要new？因为这里使用db.select，他会创建一个类似的new操作，就不需要手动new了
	sqlStr := "select community_id,community_name from community"
	//db.Select 用于查询多行并填充切片
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// 根据id查询社区分类的详情，包括社区id，社区名字，社区简介，社区的创建时间
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	//为什么这里需要new？如果不new，直接用community，默认是nil，也就是空指针，这样用db.get会报错
	community = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id = ?`
	//db.Get 用于查询单行并填充单个结构体
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}
