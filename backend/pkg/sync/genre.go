package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
)

func Genre() (err error) {

	url := "https://api.themoviedb.org/3/genre/movie/list?language=zh"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	token, err := config.GetTMDBToken()
	if err != nil {
		return fmt.Errorf("获取TMDB Token失败: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API返回错误状态码: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	type genreResponse struct {
		Genres []models.Genre `json:"genres"`
	}

	var response genreResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	for _, genre := range response.Genres {

		result := config.DB.FirstOrCreate(&genre)

		if result.Error != nil {
			return result.Error
		}
	}

	return
}
