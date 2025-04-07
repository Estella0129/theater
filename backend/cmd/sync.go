/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/pkg/sync"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var manual *bool
var interval *int

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析命令行参数

		// 初始化数据库
		config.InitDB()

		_ = sync.Genre()

		if *manual {
			// 手动执行同步
			if err := sync.SyncMovies(); err != nil {
				log.Fatalf("同步失败: %v", err)
			}
			log.Println("同步成功")
			return
		}

		// 定时同步
		duration := time.Duration(*interval) * time.Minute
		log.Printf("启动定时同步服务，间隔 %v", duration)
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := sync.SyncMovies(); err != nil {
					log.Printf("同步失败: %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	manual = syncCmd.Flags().BoolP("manual", "m", true, "手动执行同步")
	interval = syncCmd.Flags().IntP("interval", "i", 60, "定时同步间隔(分钟)")
}
