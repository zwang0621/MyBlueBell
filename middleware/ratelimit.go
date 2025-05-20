package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// gin.HandlerFunc 就是 func(*gin.Context) 的别名
// 在 RateLimitMiddleware 中，内部函数就“捕获”了在外部函数中创建的那个 bucket 变量
// 这意味着，无论有多少个请求经过这个中间件，每次执行内部函数时，它访问到的都是同一个 bucket 实例
// 这对于限流器来说是至关重要的，因为它需要追踪所有请求的状态
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//取到令牌就放行
		c.Next()
	}
}
