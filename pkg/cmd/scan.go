package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/seaung/nox/pkg/port"
	"github.com/spf13/cobra"
)

var (
	scanPorts     string
	scanTimeout   int
	scanConcurrent int
)

var scanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "端口扫描模块",
	Long:  "扫描目标主机的开放端口，支持设置端口范围和并发数量",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// 创建端口扫描实例
		ps := port.NewPortScanner(target, port.TCP_CONNECT)

		// 解析端口范围
		if strings.Contains(scanPorts, "-") {
			parts := strings.Split(scanPorts, "-")
			if len(parts) == 2 {
				start, err1 := strconv.Atoi(parts[0])
				end, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					ps.SetPortRange(start, end)
				}
			}
		} else {
			ports := make([]int, 0)
			for _, p := range strings.Split(scanPorts, ",") {
				if port, err := strconv.Atoi(p); err == nil {
					ports = append(ports, port)
				}
			}
			ps.SetPorts(ports)
		}

		// 设置超时时间和并发数
		ps.Timeout = time.Duration(scanTimeout) * time.Second
		ps.Concurrent = scanConcurrent

		// 执行端口扫描
		results := ps.Scan()

		// 输出扫描结果
		fmt.Printf("\n目标主机: %s\n", target)
		for _, result := range results {
			fmt.Printf("端口 %d: %s (%s)\n", result.Port, result.Service, result.State)
		}
		fmt.Printf("\n总计发现 %d 个开放端口\n", len(results))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// 添加命令行参数
	scanCmd.Flags().StringVarP(&scanPorts, "ports", "p", "1-1000", "端口范围 (例如: 80,443 或 1-1000)")
	scanCmd.Flags().IntVarP(&scanTimeout, "timeout", "t", 2, "单个端口扫描超时时间 (秒) (默认: 2)")
	scanCmd.Flags().IntVarP(&scanConcurrent, "concurrent", "c", 100, "并发数量 (默认: 100)")
}