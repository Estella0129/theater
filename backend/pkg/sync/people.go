package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
)

// Cast 结构体对应TMDB API返回的演员数据
type Cast struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CreditID           string  `json:"credit_id"`
	CastID             int     `json:"cast_id"`   // 演员ID
	Character          string  `json:"character"` // 饰演角色
	Order              int     `json:"order"`     // 演员排序
}

// Crew 结构体对应TMDB API返回的演职人员数据
type Crew struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CreditID           string  `json:"credit_id"`
	Department         string  `json:"department"`
	Job                string  `json:"job"`
}

type PeoplesResponse struct {
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}

func SyncPeople(movieID int) (err error) {

	var body []byte

	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?language=zh-CN", movieID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	token, err := config.GetTMDBToken()
	if err != nil {
		return fmt.Errorf("获取TMDB Token失败: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ = io.ReadAll(res.Body)

	var data PeoplesResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("解析JSON失败: %v", err)
	}

	for index, item := range data.Crew {

		err = syncCredit(movieID, item.CreditID, index)
		if err != nil {
			fmt.Println(err)
		}
	}

	for index, item := range data.Cast {

		err = syncCredit(movieID, item.CreditID, index)
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

type CreditResponse struct {
	CreditType string `json:"credit_type"`
	Department string `json:"department"`
	Job        string `json:"job"`
	Media      struct {
		BackdropPath     string  `json:"backdrop_path"`
		ID               int     `json:"id"`
		Title            string  `json:"title"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path"`
		MediaType        string  `json:"media_type"`
		Adult            bool    `json:"adult"`
		OriginalLanguage string  `json:"original_language"`
		GenreIds         []int   `json:"genre_ids"`
		Popularity       float64 `json:"popularity"`
		ReleaseDate      string  `json:"release_date"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
		Character        string  `json:"character"`
	} `json:"media"`
	MediaType string `json:"media_type"`
	ID        string `json:"id"`
	People    struct {
		ID                 int         `json:"id"`
		Name               string      `json:"name"`
		OriginalName       string      `json:"original_name"`
		MediaType          string      `json:"media_type"`
		Adult              bool        `json:"adult"`
		Popularity         float64     `json:"popularity"`
		Gender             int         `json:"gender"`
		KnownForDepartment string      `json:"known_for_department"`
		ProfilePath        interface{} `json:"profile_path"`
	} `json:"People"`
}

func syncCredit(movieID int, id string, index int) (err error) {
	var dbCredit models.Credit
	config.DB.Where("credit_id = ?", id).First(&dbCredit)

	if dbCredit.ID != "" {
		return
	}
	fmt.Println(fmt.Printf("syncCredit: %s\n", id))

	var response *CreditResponse
	response, err = getCreditResponse(id)
	if err != nil {
		return
	}

	err = syncPeople(response.People.ID)
	if err != nil {
		return
	}

	dbCredit = models.Credit{
		ID:         response.ID,
		MovieID:    movieID,
		CreditType: response.CreditType,
		Department: response.Department,
		Job:        response.Job,
		PeopleID:   response.People.ID,
		Order:      index,
	}

	err = config.DB.Create(&dbCredit).Error
	return
}

func syncPeople(id int) (err error) {

	var People models.People
	config.DB.Where("id = ?", id).First(&People)

	if People.ID != 0 {
		return
	}
	fmt.Println(fmt.Printf("syncPeople: %d\n", id))

	url := fmt.Sprintf("https://api.themoviedb.org/3/People/%d?language=zh-CN", id)

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

	err = json.Unmarshal(body, &People)

	err = config.DB.Create(&People).Error

	return
}

func getCreditResponse(id string) (response *CreditResponse, err error) {

	url := fmt.Sprintf("https://api.themoviedb.org/3/credit/%s", id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	token, err := config.GetTMDBToken()
	if err != nil {
		return nil, fmt.Errorf("获取TMDB Token失败: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	response = &CreditResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return
}
