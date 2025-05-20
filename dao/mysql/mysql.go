package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	//DSN:Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	//也可以使用mustconnect连接不成功就panic
	//相当于sql的connect加ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
	}
	//数值需要根据业务情况具体确定
	db.SetConnMaxIdleTime(time.Second * 10)
	db.SetConnMaxLifetime(time.Second * 100)
	db.SetMaxIdleConns(viper.GetInt("mysql.MaxOpenConns"))
	db.SetMaxOpenConns(viper.GetInt("mysql.MaxIdleConns"))

	return nil
}

func Close() {
	_ = db.Close()
}
