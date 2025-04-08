/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/handlers"
	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库连接
		config.InitDB()

		// 创建Gin路由引擎
		r := gin.Default()

		// 设置API路由
		v1 := r.Group("/api/v1")
		{
			// 前端接口路由组
			frontend := v1.Group("/frontend")
			{
				// 用户相关路由
				frontend.POST("/users/register", handlers.RegisterUser) // 用户注册
				frontend.POST("/users/login", handlers.LoginUser)       // 用户登录
				frontend.GET("/users/:id", handlers.GetUser)            // 获取用户详情

				// 电影相关路由
				frontend.GET("/movies", handlers.GetMovies)    // 获取电影列表
				frontend.GET("/movies/:id", handlers.GetMovie) // 获取单个电影详情
				frontend.GET("/genres", handlers.GetGenres)    // 获取所有电影类型
			}

			// 管理后台接口路由组
			admin := v1.Group("/admin")
			{
				admin.Use()
				// 用户管理路由
				admin.POST("/users", handlers.CreateUser)       // 管理员创建用户
				admin.GET("/users", handlers.GetUsers)          // 获取用户列表
				admin.PUT("/users/:id", handlers.UpdateUser)    // 更新用户信息
				admin.DELETE("/users/:id", handlers.DeleteUser) // 删除用户

				// 电影管理路由
				admin.POST("/movies", handlers.CreateMovie)       // 创建电影
				admin.GET("/movies", handlers.GetAdminMovies)     // 获取电影列表
				admin.PUT("/movies/:id", handlers.UpdateMovie)    // 更新电影信息
				admin.DELETE("/movies/:id", handlers.DeleteMovie) // 删除电影
			}
		}

		// 启动HTTP服务器
		r.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
