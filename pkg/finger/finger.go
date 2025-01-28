package finger

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/projectdiscovery/wappalyzergo"
	"github.com/seaung/nox/pkg/utils"
)

// FingerScanner 网站指纹识别扫描器结构体
type FingerScanner struct {
	Target   string         // 目标URL
	Timeout  time.Duration // 请求超时时间
	Logger   *utils.Logger // 日志记录器
}

// FingerResult 指纹识别结果结构体
type FingerResult struct {
	URL         string   // 目标URL
	Technologies []string // 识别到的技术列表
}

// NewFingerScanner 创建一个新的指纹识别扫描器实例
func NewFingerScanner(target string) *FingerScanner {
	return &FingerScanner{
		Target:   target,
		Timeout:  time.Second * 10,
		Logger:   utils.New(),
	}
}

// SetTimeout 设置请求超时时间
func (fs *FingerScanner) SetTimeout(timeout time.Duration) {
	fs.Timeout = timeout
}

// Scan 执行指纹识别扫描
func (fs *FingerScanner) Scan() (*FingerResult, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: fs.Timeout,
	}

	// 发送HTTP请求
	resp, err := client.Get(fs.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 创建wappalyzer实例
	wappalyzer, err := wappalyzergo.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create wappalyzer instance: %v", err)
	}

	// 识别技术
	techs := wappalyzer.Fingerprint(resp.Header, body)
	technologies := make([]string, 0)
	for tech := range techs {
		technologies = append(technologies, tech)
		fs.Logger.Success(fmt.Sprintf("Detected technology: %s", tech))
	}

	return &FingerResult{
		URL:         fs.Target,
		Technologies: technologies,
	}, nil
}