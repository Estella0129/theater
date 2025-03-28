package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
)

// SyncMovies 从TiDB同步电影信息并写入本地数据库
func SyncMovies() error {
	// 1. 从TMDB API获取电影数据
	var err error
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=zh-CN&page=2&sort_by=popularity.desc"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI1M2EyN2MzYzA1MDQ0Mzk0ZGE1NjQ0NTdhYmVlNWY1ZCIsIm5iZiI6MTcyNzQxOTQwMS42MjcsInN1YiI6IjY2ZjY1NDA5YWE3ZTVmYTIwMjk2NWE5ZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.miOSW9RORTxh-vNgGgZhmIjkWqrZX3TIBreHJJOZKVQ")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var tmdbResponse struct {
		Results []struct {
			ID               int     `json:"id"`
			Title            string  `json:"title"`
			OriginalTitle    string  `json:"original_title"`
			OriginalLanguage string  `json:"original_language"`
			Overview         string  `json:"overview"`
			PosterPath       string  `json:"poster_path"`
			BackdropPath     string  `json:"backdrop_path"`
			ReleaseDate      string  `json:"release_date"`
			Adult            bool    `json:"adult"`
			Popularity       float64 `json:"popularity"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
			Video            bool    `json:"video"`
			GenreIDs         []int   `json:"genre_ids"`
		}
	}

	if err := json.Unmarshal(body, &tmdbResponse); err != nil {
		return err
	}

	// 2. 数据转换和处理
	for _, tmdbMovie := range tmdbResponse.Results {

		releaseDate, _ := time.Parse("2006-01-02", tmdbMovie.ReleaseDate)

		movie := models.Movie{
			ID:               uint(tmdbMovie.ID),
			Title:            tmdbMovie.Title,
			OriginalTitle:    tmdbMovie.OriginalTitle,
			OriginalLanguage: tmdbMovie.OriginalLanguage,
			Overview:         tmdbMovie.Overview,
			PosterPath:       tmdbMovie.PosterPath,
			BackdropPath:     tmdbMovie.BackdropPath,
			ReleaseDate:      releaseDate,
			Adult:            tmdbMovie.Adult,
			Popularity:       tmdbMovie.Popularity,
			VoteAverage:      tmdbMovie.VoteAverage,
			VoteCount:        tmdbMovie.VoteCount,
			Video:            tmdbMovie.Video,
		}

		// 3. 使用GORM保存到SQLite
		result := config.DB.Create(&movie)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// StartSyncService 启动定时同步服务
func main() {
	// 解析命令行参数
	manual := flag.Bool("manual", false, "手动执行同步")
	interval := flag.Int("interval", 60, "定时同步间隔(分钟)")
	flag.Parse()

	// 初始化数据库
	config.InitDB()

	if *manual {
		// 手动执行同步
		if err := SyncMovies(); err != nil {
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
			if err := SyncMovies(); err != nil {
				log.Printf("同步失败: %v", err)
			}
		}
	}
}

func StartSyncService() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时同步一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := SyncMovies(); err != nil {
				log.Printf("同步电影信息失败: %v", err)
			}
		}
	}
}

// SetupSyncRoutes 设置同步相关的API路由
func SetupSyncRoutes(r *gin.Engine) {
	r.POST("/api/sync/movies", func(c *gin.Context) {
		if err := SyncMovies(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "电影同步成功"})
	})
}
