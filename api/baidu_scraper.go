// file: internal/api/baidu_scraper.go (或者你喜欢的任何位置)

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"web_app/models" // 替换成你的项目路径

	"github.com/PuerkitoBio/goquery"
)

// GetBaiduRealtimeHotNews 抓取并解析百度实时热点
func GetBaiduRealtimeHotNews() ([]*models.NewsItem, error) {
	// 1. 目标 URL
	url := "https://top.baidu.com/board?tab=realtime"

	// 2. 创建 HTTP 请求
	// 设置请求头，模拟浏览器访问，否则可能被目标网站拒绝
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 4. 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 5. 提取数据
	var newsList []*models.NewsItem
	// 通过分析百度热榜网页的 HTML 结构，我们找到包含所有热点的父容器
	// 注意：这些 class 名称可能会随百度网页改版而改变
	doc.Find("div.category-wrap_iQLoo").Each(func(i int, s *goquery.Selection) {
		// 提取标题和链接
		titleSelector := s.Find("div.content_1YWBm a")
		title := titleSelector.Text()
		url, _ := titleSelector.Attr("href")

		// 提取热度
		heatSelector := s.Find("div.hot-index_1Bl1a")
		heat := heatSelector.Text()

		// 提取排名
		rankSelector := s.Find("div.index_1Ew5p")
		rank, _ := strconv.Atoi(rankSelector.Text())

		if title != "" && url != "" {
			newsList = append(newsList, &models.NewsItem{
				Rank:  rank,
				Title: strings.TrimSpace(title),
				URL:   url,
				Heat:  strings.TrimSpace(heat),
			})
		}
	})

	if len(newsList) == 0 {
		return nil, fmt.Errorf("未能从页面提取到任何新闻条目，可能是网页结构已更改")
	}

	return newsList, nil
}
