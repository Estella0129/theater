package sync

import (
	"encoding/json"
	"fmt"
	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"io"
	"net/http"
	"time"
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
