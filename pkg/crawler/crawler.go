package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/seaung/nox/pkg/utils"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

// Crawler Web爬虫结构体
type Crawler struct {
	Target     string         // 目标URL
	Depth      int           // 爬取深度
	Timeout    time.Duration // 超时时间
	Concurrent int          // 并发数量
	Logger     *utils.Logger // 日志记录器

	visited sync.Map        // 已访问的URL
	browser *rod.Browser    // rod浏览器实例
}

// CrawlResult 爬取结果结构体
type CrawlResult struct {
	URL      string   // 发现的URL
	Depth    int     // URL的深度
	ParentURL string // 父URL
}

// NewCrawler 创建一个新的爬虫实例
func NewCrawler(target string) *Crawler {
	return &Crawler{
		Target:     target,
		Depth:      3,
		Timeout:    time.Second * 30,
		Concurrent: 5,
		Logger:     utils.New(),
	}
}

// SetDepth 设置爬取深度
func (c *Crawler) SetDepth(depth int) {
	c.Depth = depth
}

// SetTimeout 设置超时时间
func (c *Crawler) SetTimeout(timeout time.Duration) {
	c.Timeout = timeout
}

// SetConcurrent 设置并发数量
func (c *Crawler) SetConcurrent(concurrent int) {
	c.Concurrent = concurrent
}

// extractURLsFromJS 从JavaScript代码中提取URL
func (c *Crawler) extractURLsFromJS(jsCode string) []string {
	urls := make([]string, 0)
	// 解析JavaScript代码
	lexer := js.NewLexer(parse.NewInput(strings.NewReader(jsCode)))

	for {
		tt, text := lexer.Next()
		if tt == js.ErrorToken {
			break
		}

		// 检查字符串中是否包含URL
		if tt == js.StringToken || tt == js.TemplateToken {
			str := string(text)
			if strings.HasPrefix(str, "http") || strings.HasPrefix(str, "/") {
				urls = append(urls, str)
			}
		}
	}

	return urls
}

// normalizeURL 规范化URL
func (c *Crawler) normalizeURL(rawURL, parentURL string) string {
	// 处理相对路径
	if strings.HasPrefix(rawURL, "/") {
		parentURLObj, err := url.Parse(parentURL)
		if err != nil {
			return ""
		}
		rawURL = fmt.Sprintf("%s://%s%s", parentURLObj.Scheme, parentURLObj.Host, rawURL)
	}

	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	// 移除URL中的片段标识符
	parsedURL.Fragment = ""

	return parsedURL.String()
}

// crawlPage 爬取单个页面
func (c *Crawler) crawlPage(pageURL string, depth int, resultsChan chan<- *CrawlResult) {
	// 检查深度限制
	if depth > c.Depth {
		return
	}

	// 检查URL是否已访问
	if _, visited := c.visited.LoadOrStore(pageURL, true); visited {
		return
	}

	// 创建新的页面
	page := c.browser.MustPage(pageURL)
	defer page.Close()

	// 设置超时
	page.Timeout(c.Timeout)

	// 等待页面加载完成
	page.MustWaitLoad()

	// 获取所有链接
	elements := page.MustElements("a[href]")
	for _, element := range elements {
		href := element.MustAttribute("href")
		if href == nil {
			continue
		}

		normalizedURL := c.normalizeURL(*href, pageURL)
		if normalizedURL != "" {
			resultsChan <- &CrawlResult{
				URL:       normalizedURL,
				Depth:     depth,
				ParentURL: pageURL,
			}
		}
	}

	// 获取并分析页面上的JavaScript代码
	scripts := page.MustElements("script")
	for _, script := range scripts {
		// 获取内联JavaScript
		if jsCode, err := script.Text(); err == nil && jsCode != "" {
			urls := c.extractURLsFromJS(jsCode)
			for _, rawURL := range urls {
				normalizedURL := c.normalizeURL(rawURL, pageURL)
				if normalizedURL != "" {
					resultsChan <- &CrawlResult{
						URL:       normalizedURL,
						Depth:     depth,
						ParentURL: pageURL,
					}
				}
			}
		}

		// 获取外部JavaScript文件
		if src := script.MustAttribute("src"); src != nil {
			normalizedURL := c.normalizeURL(*src, pageURL)
			if normalizedURL != "" {
				resultsChan <- &CrawlResult{
					URL:       normalizedURL,
					Depth:     depth,
					ParentURL: pageURL,
				}
			}
		}
	}
}

// Crawl 执行爬虫任务
func (c *Crawler) Crawl() []*CrawlResult {
	// 启动浏览器
	url := launcher.New().
		Headless(true).
		Leakless(true).
		MustLaunch()

	c.browser = rod.New().ControlURL(url).MustConnect()
	defer c.browser.MustClose()

	results := make([]*CrawlResult, 0)
	resultsChan := make(chan *CrawlResult, 1000)
	wg := sync.WaitGroup{}

	// 创建工作池
	jobs := make(chan string, c.Concurrent)
	for i := 0; i < c.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				c.crawlPage(url, 1, resultsChan)
			}
		}()
	}

	// 发送初始URL
	jobs <- c.Target

	// 处理结果
	go func() {
		for result := range resultsChan {
			c.Logger.Success(fmt.Sprintf("Found URL: %s (Depth: %d, Parent: %s)", result.URL, result.Depth, result.ParentURL))
			results = append(results, result)

			// 如果深度未达到限制，将新URL加入任务队列
			if result.Depth < c.Depth {
				jobs <- result.URL
			}
		}
	}()

	// 等待所有任务完成
	wg.Wait()
	close(jobs)
	close(resultsChan)

	return results
}