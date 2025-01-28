package port

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/seaung/nox/pkg/utils"
)

// ScanType 定义端口扫描的类型
type ScanType int

const (
	// TCP_CONNECT 使用TCP连接方式扫描
	TCP_CONNECT ScanType = iota
	// TCP_SYN 使用TCP SYN方式扫描
	TCP_SYN
	// UDP 使用UDP方式扫描
	UDP
)

// PortScanner 端口扫描器结构体
type PortScanner struct {
	Target     string         // 目标主机地址
	Ports      []int         // 要扫描的端口列表
	Timeout    time.Duration // 连接超时时间
	Concurrent int          // 并发扫描的数量
	ScanType   ScanType      // 扫描类型
	Logger     *utils.Logger // 日志记录器
}

// ScanResult 端口扫描结果结构体
type ScanResult struct {
	Port    int    // 端口号
	State   string // 端口状态（open/closed/filtered）
	Service string // 端口对应的服务名称
}

// NewPortScanner 创建一个新的端口扫描器实例
func NewPortScanner(target string, scanType ScanType) *PortScanner {
	return &PortScanner{
		Target:     target,
		Timeout:    time.Second * 2,
		Concurrent: 100,
		ScanType:   scanType,
		Logger:     utils.New(),
	}
}

// SetPorts 设置要扫描的端口列表
func (ps *PortScanner) SetPorts(ports []int) {
	ps.Ports = ports
}

// SetPortRange 设置要扫描的端口范围
func (ps *PortScanner) SetPortRange(start, end int) error {
	if start > end || start < 1 || end > 65535 {
		return fmt.Errorf("invalid port range")
	}
	ports := make([]int, 0)
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	ps.Ports = ports
	return nil
}

// tcpConnect 使用TCP连接方式扫描单个端口
func (ps *PortScanner) tcpConnect(port int) ScanResult {
	target := fmt.Sprintf("%s:%d", ps.Target, port)
	conn, err := net.DialTimeout("tcp", target, ps.Timeout)
	result := ScanResult{Port: port}
	
	if err != nil {
		result.State = "closed"
		return result
	}
	defer conn.Close()
	
	result.State = "open"
	result.Service = getServiceName(port)
	return result
}

// udpScan 使用UDP方式扫描单个端口
func (ps *PortScanner) udpScan(port int) ScanResult {
	target := fmt.Sprintf("%s:%d", ps.Target, port)
	conn, err := net.DialTimeout("udp", target, ps.Timeout)
	result := ScanResult{Port: port}
	
	if err != nil {
		result.State = "closed"
		return result
	}
	defer conn.Close()
	
	// UDP需要发送数据才能确定端口状态
	_, err = conn.Write([]byte("\x00"))
	if err != nil {
		result.State = "closed"
		return result
	}
	
	result.State = "open|filtered"
	result.Service = getServiceName(port)
	return result
}

// Scan 执行端口扫描
// 使用goroutine实现并发扫描，通过channel进行任务分发和结果收集
func (ps *PortScanner) Scan() []ScanResult {
	results := make([]ScanResult, 0)
	jobs := make(chan int, len(ps.Ports))
	resultsChan := make(chan ScanResult, len(ps.Ports))
	wg := sync.WaitGroup{}

	// 启动工作协程
	for i := 0; i < ps.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				var result ScanResult
				switch ps.ScanType {
				case TCP_CONNECT:
					result = ps.tcpConnect(port)
				case UDP:
					result = ps.udpScan(port)
				default:
					result = ps.tcpConnect(port)
				}
				resultsChan <- result
			}
		}()
	}

	// 发送任务
	go func() {
		for _, port := range ps.Ports {
			jobs <- port
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
		if result.State == "open" || result.State == "open|filtered" {
			ps.Logger.Success(fmt.Sprintf("Port %d is %s (%s)", result.Port, result.State, result.Service))
			results = append(results, result)
		}
	}

	return results
}

// getServiceName 根据端口号获取对应的服务名称
func getServiceName(port int) string {
	commonPorts := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		3306: "MySQL",
		5432: "PostgreSQL",
		6379: "Redis",
		8080: "HTTP-Proxy",
	}

	if service, ok := commonPorts[port]; ok {
		return service
	}
	return "unknown"
}