package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	fmt.Println(err)
	// 分页获取所有电影数据
	page := 1
	totalPages := 1
	var allResults []struct {
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

	for page <= totalPages {
		// 添加适当的延迟避免API限流
		time.Sleep(500 * time.Millisecond)

		url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=zh-CN&page=%d&sort_by=popularity.desc", page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}

		req.Header.Add("accept", "application/json")
		token, err := config.GetTMDBToken()
		if err != nil {
			return fmt.Errorf("获取TMDB Token失败: %v", err)
		}
		req.Header.Add("Authorization", "Bearer "+token)

		fmt.Println("request", url)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("请求失败: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("API返回错误状态码: %d", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("读取响应体失败: %v", err)
		}

		var tmdbResponse struct {
			Page         int `json:"page"`
			TotalPages   int `json:"total_pages"`
			TotalResults int `json:"total_results"`
			Results      []struct {
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
			return fmt.Errorf("解析JSON失败: %v", err)
		}

		totalPages = tmdbResponse.TotalPages
		if totalPages > 10 {
			totalPages = 10
		}
		allResults = append(allResults, tmdbResponse.Results...)
		page++
	}

	// 使用收集到的所有结果进行处理
	for _, tmdbMovie := range allResults {

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
		result := config.DB.FirstOrCreate(&movie)
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

// GetMovies 从本地数据库获取电影信息
func GetMovies(c *gin.Context) {
	var movies []models.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, movies)
}
