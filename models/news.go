package models

// NewsItem 定义单条新闻的数据结构
type NewsItem struct {
	Rank  int    `json:"rank"`  // 排名
	Title string `json:"title"` // 标题
	URL   string `json:"url"`   // 新闻链接
	Heat  string `json:"heat"`  // 热度值
}

// ApiNewsResponse 定义最终返回给前端的完整结构
type ApiNewsResponse struct {
	Total int64       `json:"total"` // 总条数
	List  []*NewsItem `json:"list"`  // 新闻列表
}
