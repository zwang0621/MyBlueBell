package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routers"
	"web_app/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go web 开发较通用的脚手架
// @title 汪壮的bluebell
// @version 1.0
// @description 这是汪壮制作的第一个golang项目blubell
// @termsOfService http://swagger.io/terms/

// @contact.name wangzhuang
// @contact.url http://www.swagger.io/support
// @contact.email wz990621@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080/api/v1

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(viper.GetString("app.mode")); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success!")
	//3.初始化mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()
	//4.初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()
	//初始化雪花算法，获取分布式id
	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("Init snowflake failed,err:%v\n", err)
		return
	}
	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed,err:%v\n", err)
		return
	}
	//5.注册路由
	r := routers.Setup(viper.GetString("app.mode"))
	r.Run()
	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:%s\n", zap.Error(err))
		}
	}()

	//创建一个接受系统信号的通道
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit //当收到上述两种信号才会往下执行
	zap.L().Info("Shutdown Server ...")
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//5秒内优雅关闭服务（将未处理完的请求处理完在关闭服务），超过五秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting ...")
}
