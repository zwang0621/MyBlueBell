package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {

	//1.判断用户存不存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}

	//2.生成uid
	userID := snowflake.GenID()

	//3.构建一个user实例
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//4.保存进数据库
	return mysql.InsertUser(&user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {

	// 判断用户存不存在
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到userid
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt token
	atoken, rtoken, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.AccessToken = atoken
	user.RefreshToken = rtoken
	return
}
