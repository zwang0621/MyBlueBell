package routers

import (
	"net/http"
	"time"
	"web_app/controller"
	"web_app/logger"
	"web_app/middleware"

	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

/*
1.GET: 用于从服务器获取数据。它应该是“安全”和“幂等”的
2.POST: 用于向服务器发送数据，通常用于创建、更新资源，或者执行有副作用的操作。它不是幂等的，也通常不是安全的

GET: 是“幂等”的。对同一个 URL 发送多次相同的 GET 请求，其结果应该是相同的（都是获取同样的数据）
POST: 通常不是幂等的。发送多次相同的 POST 请求，可能会导致不同的结果（例如，多次提交同一个评论表单，可能会创建多条重复评论）
*/
func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//将限速中间件应用到 r.Use() 上是合理的，因为它需要对所有进入你应用层的请求都生效，这也包括集成zap日志库到gin中
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleware.RateLimitMiddleware(time.Second*2, 1))

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	//登陆业务路由
	v1.POST("/login", controller.LoginHandler)
	//刷新atoken路由
	v1.POST("/refreshtoken", controller.RefreshTokenHandler)

	// 搜索业务-搜索帖子
	v1.GET("/search", controller.PostSearchHandler)
	//查看帖子列表其实不需要登陆
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	//查看社区信息也不需要登陆
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	//增加侧边栏热点新闻功能
	v1.GET("/news", controller.NewsTrendingHandler)

	//创建一个v1路由组的好处：
	//应用中间件 (Applying Middleware): 可以方便地将一组中间件应用到一组相关的路由上。而不是对每个路由都单独调用 Use()
	//如果将来你需要开发 /api/v2 版本的 API，你可以很容易地创建另一个组 v2 := r.Group("/api/v2")，并在其中定义新的路由和应用不同的中间件，而不会影响到现有的 v1 组的代码
	v1.Use(middleware.JWTAuthMiddleware()) //应用jwt认证中间件

	{

		v1.POST("/post", controller.CreatePostHandler)

		v1.POST("/vote", controller.PostVoteController)
	}

	r.POST("/ping", middleware.JWTAuthMiddleware(), func(ctx *gin.Context) {
		//如果是登录的用户（如何判断？判断请求头是否有有效的jwt token,已经用中间件实现）
		ctx.String(http.StatusOK, "pong")
	})

	pprof.Register(r) //注册pprof相关路由
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
