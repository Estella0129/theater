package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
)

func Images(movieID int) (err error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/images", movieID)

	fmt.Println("request: ", url)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	token, err := config.GetTMDBToken()
	if err != nil {
		return fmt.Errorf("获取TMDB Token失败: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	if res != nil {
		defer res.Body.Close()
	}
	if res == nil || res.StatusCode != http.StatusOK {
		if res != nil {
			return fmt.Errorf("请求失败，状态码: %d", res.StatusCode)
		}
		return fmt.Errorf("请求未成功，响应为 nil")
	}
	body, _ := io.ReadAll(res.Body)

	type imagesResponse struct {
		Backdrops []models.Image `json:"backdrops"`
		Posters   []models.Image `json:"posters"`
		Logos     []models.Image `json:"logos"`
	}

	var response imagesResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	for index, item := range response.Backdrops {

		if item.Iso6391 != "" && item.Iso6391 != "zh" {
			continue
		}

		r := models.MovieImage{
			MovieID:       movieID,
			ImageFilePath: item.FilePath,
		}
		result := config.DB.FirstOrCreate(&r)
		if result.Error != nil {
			return result.Error
		}

		item.Type = "backdrop"
		result = config.DB.FirstOrCreate(&item)

		if result.Error != nil {
			return result.Error
		}

		if index > 3 {
			break
		}
	}

	for index, item := range response.Posters {

		if item.Iso6391 != "" && item.Iso6391 != "zh" {
			continue
		}
		r := models.MovieImage{
			MovieID:       movieID,
			ImageFilePath: item.FilePath,
		}
		result := config.DB.FirstOrCreate(&r)
		if result.Error != nil {
			return result.Error
		}

		item.Type = "poster"
		result = config.DB.FirstOrCreate(&item)

		if result.Error != nil {
			return result.Error
		}

		if index > 3 {
			break
		}
	}

	for index, item := range response.Logos {

		if item.Iso6391 != "" && item.Iso6391 != "zh" {
			continue
		}

		r := models.MovieImage{
			MovieID:       movieID,
			ImageFilePath: item.FilePath,
		}
		result := config.DB.FirstOrCreate(&r)
		if result.Error != nil {
			return result.Error
		}

		item.Type = "logo"
		result = config.DB.FirstOrCreate(&item)

		if result.Error != nil {
			return result.Error
		}

		if index > 3 {
			break
		}
	}

	return nil
}
