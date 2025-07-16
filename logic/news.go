package logic

import (
	"web_app/api"
	"web_app/models"
)

// GetBaiduNewsTrending 获取百度新闻列表 (函数名建议修改)
func GetBaiduNewsTrending(p *models.ParamNewsTrending) (*models.ApiNewsResponse, error) {
	// 调用我们新的抓取函数来获取所有新闻数据
	allNews, err := api.GetBaiduRealtimeHotNews()
	if err != nil {
		return nil, err // 如果抓取失败，直接返回错误
	}

	// --- 处理分页逻辑 ---
	total := int64(len(allNews))

	// 计算分页的起始和结束索引
	startIndex := (p.Offset - 1) * p.Limit
	if startIndex < 0 || startIndex >= total {
		// 如果请求的页码超出范围，返回一个空列表
		return &models.ApiNewsResponse{
			Total: total,
			List:  []*models.NewsItem{},
		}, nil
	}

	endIndex := startIndex + p.Limit
	if endIndex > total {
		endIndex = total // 防止索引越界
	}

	// 从所有新闻中切片出当前页的数据
	OffsetdNews := allNews[startIndex:endIndex]

	// 组装最终返回的数据结构
	res := &models.ApiNewsResponse{
		Total: total,
		List:  OffsetdNews,
	}

	return res, nil
}
