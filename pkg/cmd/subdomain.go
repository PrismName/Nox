package cmd

import (
	"fmt"
	"time"

	"github.com/seaung/nox/pkg/subdomain"
	"github.com/spf13/cobra"
)

var (
	subdomainWordlist string
	subdomainTimeout  int
	subdomainConcurrent int
)

var subdomainCmd = &cobra.Command{
	Use:   "subdomain [domain]",
	Short: "子域名扫描模块",
	Long:  "扫描目标域名的子域名，支持自定义字典和并发数量",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// 创建子域名扫描实例
		ss := subdomain.NewSubdomainScanner(target)
		ss.SetWordlist([]string{subdomainWordlist})
		ss.SetTimeout(time.Duration(subdomainTimeout) * time.Second)
		ss.SetConcurrent(subdomainConcurrent)

		// 执行子域名扫描
		results := ss.Scan()

		// 输出扫描结果
		fmt.Printf("\n目标域名: %s\n", target)
		for _, result := range results {
			fmt.Printf("子域名: %s (IP: %v)\n", result.Subdomain, result.IPList)
		}
		fmt.Printf("\n总计发现 %d 个子域名\n", len(results))
	},
}

func init() {
	rootCmd.AddCommand(subdomainCmd)

	// 添加命令行参数
	subdomainCmd.Flags().StringVarP(&subdomainWordlist, "wordlist", "w", "wordlist/domain.nox", "子域名字典文件路径")
	subdomainCmd.Flags().IntVarP(&subdomainTimeout, "timeout", "t", 5, "单个子域名解析超时时间 (秒) (默认: 5)")
	subdomainCmd.Flags().IntVarP(&subdomainConcurrent, "concurrent", "c", 50, "并发数量 (默认: 50)")
}