package sync

import (
	"encoding/json"
	"fmt"
	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"io"
	"net/http"
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

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

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
