package dirs

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/seaung/nox/pkg/utils"
)

// DirScanner 目录扫描器结构体
type DirScanner struct {
	Target     string         // 目标URL
	Wordlist   []string      // 字典列表
	Timeout    time.Duration // HTTP请求超时时间
	Concurrent int          // 并发数量
	Logger     *utils.Logger // 日志记录器
}

// DirResult 目录扫描结果结构体
type DirResult struct {
	Path       string // 目录路径
	StatusCode int    // HTTP状态码
	Length     int64  // 响应内容长度
}

// NewDirScanner 创建一个新的目录扫描器实例
func NewDirScanner(target string) *DirScanner {
	return &DirScanner{
		Target:     target,
		Timeout:    time.Second * 10,
		Concurrent: 50,
		Logger:     utils.New(),
	}
}

// SetWordlist 设置字典列表
func (ds *DirScanner) SetWordlist(wordlist []string) {
	ds.Wordlist = wordlist
}

// SetTimeout 设置HTTP请求超时时间
func (ds *DirScanner) SetTimeout(timeout time.Duration) {
	ds.Timeout = timeout
}

// SetConcurrent 设置并发数量
func (ds *DirScanner) SetConcurrent(concurrent int) {
	ds.Concurrent = concurrent
}

// checkDir 检查单个目录
func (ds *DirScanner) checkDir(path string) *DirResult {
	// 构造完整URL
	url := ds.Target
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += path

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: ds.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不跟随重定向
		},
	}

	// 发送HTTP请求
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// 判断目录是否有效
	if resp.StatusCode != http.StatusNotFound {
		return &DirResult{
			Path:       path,
			StatusCode: resp.StatusCode,
			Length:     int64(len(body)),
		}
	}

	return nil
}

// Scan 执行目录扫描
func (ds *DirScanner) Scan() []DirResult {
	results := make([]DirResult, 0)
	jobs := make(chan string, len(ds.Wordlist))
	resultsChan := make(chan *DirResult, len(ds.Wordlist))
	wg := sync.WaitGroup{}

	// 启动工作协程
	for i := 0; i < ds.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				result := ds.checkDir(path)
				if result != nil {
					resultsChan <- result
				}
			}
		}()
	}

	// 发送任务
	go func() {
		for _, word := range ds.Wordlist {
			jobs <- word
		}
		close(jobs)
	}()

	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果
	for result := range resultsChan {
		ds.Logger.Success(fmt.Sprintf("Found directory: %s (Status: %d, Length: %d)", result.Path, result.StatusCode, result.Length))
		results = append(results, *result)
	}

	return results
}