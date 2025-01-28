package cmd

import (
	"fmt"
	"time"

	"github.com/seaung/nox/pkg/finger"
	"github.com/spf13/cobra"
)

var (
	fingerTimeout int
)

var fingerCmd = &cobra.Command{
	Use:   "finger [url]",
	Short: "网站指纹识别模块",
	Long:  "识别目标网站使用的技术栈和框架",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// 创建指纹识别实例
		fs := finger.NewFingerScanner(target)
		fs.SetTimeout(time.Duration(fingerTimeout) * time.Second)

		// 执行指纹识别
		result, err := fs.Scan()
		if err != nil {
			fmt.Printf("指纹识别失败: %v\n", err)
			return
		}

		// 输出识别结果
		fmt.Printf("\n目标: %s\n", result.URL)
		fmt.Printf("发现技术栈: %v\n", result.Technologies)
	},
}

func init() {
	rootCmd.AddCommand(fingerCmd)

	// 添加命令行参数
	fingerCmd.Flags().IntVarP(&fingerTimeout, "timeout", "t", 10, "请求超时时间 (秒) (默认: 10)")
}