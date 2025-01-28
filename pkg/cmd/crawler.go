package cmd

import (
	"fmt"
	"time"

	"github.com/seaung/nox/pkg/crawler"
	"github.com/spf13/cobra"
)

var (
	crawlerDepth      int
	crawlerTimeout    int
	crawlerConcurrent int
)

var crawlerCmd = &cobra.Command{
	Use:   "crawler [url]",
	Short: "网站爬虫模块",
	Long:  "爬取目标网站的URL链接，支持设置爬取深度和并发数量",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// 创建爬虫实例
		c := crawler.NewCrawler(target)
		c.SetDepth(crawlerDepth)
		c.SetTimeout(time.Duration(crawlerTimeout) * time.Second)
		c.SetConcurrent(crawlerConcurrent)

		// 执行爬虫任务
		results := c.Crawl()

		// 输出统计信息
		fmt.Printf("\n总计发现 %d 个URL\n", len(results))
	},
}

func init() {
	rootCmd.AddCommand(crawlerCmd)

	// 添加命令行参数
	crawlerCmd.Flags().IntVarP(&crawlerDepth, "depth", "d", 3, "爬取深度 (默认: 3)")
	crawlerCmd.Flags().IntVarP(&crawlerTimeout, "timeout", "t", 30, "请求超时时间 (秒) (默认: 30)")
	crawlerCmd.Flags().IntVarP(&crawlerConcurrent, "concurrent", "c", 5, "并发数量 (默认: 5)")
}