package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/post", CreatePostHandler)

	body := `{
		"community_id":1,
		"title":"test",
		"content":"just a test"
	}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/post", bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//判断响应的内容是不是按预期返回了对应的状态码错误，因为在这里测试没有携带jwt token。
	assert.Contains(t, w.Body.String(), "需要登陆")

}
