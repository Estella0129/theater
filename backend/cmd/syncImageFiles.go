/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// downloadImage 下载图片到本地
func downloadImage(url string, filepath string) error {
	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP请求失败: %s", resp.Status)
	}

	// 复制响应体到文件
	_, err = io.Copy(out, resp.Body)
	return err
}

// syncImageFilesCmd represents the syncImageFiles command
var syncImageFilesCmd = &cobra.Command{
	Use:   "syncImageFiles",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitDB()

		var movieImages []models.MovieImage
		config.DB.Limit(20).Find(&movieImages)

		baseUrl := "https://image.tmdb.org/t/p/original"
		// 图片存储目录
		imageDir := "images"

		// 创建图片目录(如果不存在)
		if _, err := os.Stat(imageDir); os.IsNotExist(err) {
			err := os.Mkdir(imageDir, 0755)
			if err != nil {
				panic("无法创建图片目录: " + err.Error())
			}
		}

		for _, image := range movieImages {
			// 拼接本地文件路径
			localPath := imageDir + "/" + image.ImageFilePath

			// 检查文件是否已存在
			if _, err := os.Stat(localPath); os.IsNotExist(err) {
				// 文件不存在，下载图片
				imageUrl := baseUrl + image.ImageFilePath
				err := downloadImage(imageUrl, localPath)
				if err != nil {
					// 记录下载失败
					println("下载失败:", imageUrl, "错误:", err.Error())
					continue
				}
				println("下载成功:", imageUrl)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncImageFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncImageFilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncImageFilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
