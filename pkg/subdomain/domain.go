package subdomain

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/seaung/nox/pkg/utils"
)

// SubdomainScanner 子域名扫描器结构体
type SubdomainScanner struct {
	Domain     string         // 目标域名
	Wordlist   []string      // 字典列表
	Timeout    time.Duration // DNS查询超时时间
	Concurrent int          // 并发数量
	Logger     *utils.Logger // 日志记录器
}

// SubdomainResult 子域名扫描结果结构体
type SubdomainResult struct {
	Subdomain string   // 子域名
	IPList    []string // IP地址列表
}

// NewSubdomainScanner 创建一个新的子域名扫描器实例
func NewSubdomainScanner(domain string) *SubdomainScanner {
	return &SubdomainScanner{
		Domain:     domain,
		Timeout:    time.Second * 2,
		Concurrent: 100,
		Logger:     utils.New(),
	}
}

// SetWordlist 设置字典列表
func (ss *SubdomainScanner) SetWordlist(wordlist []string) {
	ss.Wordlist = wordlist
}

// SetTimeout 设置DNS查询超时时间
func (ss *SubdomainScanner) SetTimeout(timeout time.Duration) {
	ss.Timeout = timeout
}

// SetConcurrent 设置并发数量
func (ss *SubdomainScanner) SetConcurrent(concurrent int) {
	ss.Concurrent = concurrent
}

// dnsLookup 执行DNS查询
func (ss *SubdomainScanner) dnsLookup(subdomain string) *SubdomainResult {
	ips, err := net.LookupHost(subdomain)
	if err != nil {
		return nil
	}

	return &SubdomainResult{
		Subdomain: subdomain,
		IPList:    ips,
	}
}

// Scan 执行子域名扫描
func (ss *SubdomainScanner) Scan() []SubdomainResult {
	results := make([]SubdomainResult, 0)
	jobs := make(chan string, len(ss.Wordlist))
	resultsChan := make(chan *SubdomainResult, len(ss.Wordlist))
	wg := sync.WaitGroup{}

	// 启动工作协程
	for i := 0; i < ss.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for subdomain := range jobs {
				fullDomain := fmt.Sprintf("%s.%s", subdomain, ss.Domain)
				result := ss.dnsLookup(fullDomain)
				if result != nil {
					resultsChan <- result
				}
			}
		}()
	}

	// 发送任务
	go func() {
		for _, word := range ss.Wordlist {
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
		ss.Logger.Success(fmt.Sprintf("Found subdomain: %s (IPs: %v)", result.Subdomain, result.IPList))
		results = append(results, *result)
	}

	return results
}