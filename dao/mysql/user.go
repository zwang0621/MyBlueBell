package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

//把每一步数据库的操作封装成函数

const secret = "wangzhuang"

// CheckUserExist 检查指定用户的用户名是否存在
func CheckUserExist(username string) (err error) {
	var count int
	sqlStr := `select count(user_id) from user where username = ?`
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	password := encryptPassword(user.Password)
	//执行sql语句入库
	sqlStr := `insert into user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, password)
	return
}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		//数据库里没有该用户
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库失败
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	//user.Password是数据库里存储的密码，password是用户输入的密码
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// 根据id获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
