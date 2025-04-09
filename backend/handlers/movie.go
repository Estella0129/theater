package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
)

// GetMovies 获取电影列表，支持分页和搜索
func GetMovies(c *gin.Context) {
	// 获取并验证分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	searchQuery := strings.TrimSpace(c.Query("query"))
	genre := strings.TrimSpace(c.Query("genre"))

	var movies []models.Movie
	var total int64

	offset := (page - 1) * pageSize

	dbQuery := config.DB.Model(&models.Movie{})
	if searchQuery != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR original_title LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}
	if genre != "" {
		dbQuery = dbQuery.Joins(
			"left JOIN movie_genres ON movies.id = movie_genres.movie_id",
			genre,
		).Where("movie_genres.genre_id = ?", genre)
	}

	// 获取总记录数
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取电影总数失败"})
		return
	}

	// 获取分页数据
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取电影列表失败"})
		return
	}

	// 计算总页数
	totalPages := int64(0)
	if total > 0 {
		totalPages = (total + int64(pageSize) - 1) / int64(pageSize)
	}

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages,
		"results":     movies,
	})
}

// GetMovie 获取单个电影详情
func GetMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	result := config.DB.Preload("Credits").First(&movie, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// CreateMovie 创建电影
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result := config.DB.Create(&movie)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie 更新电影信息
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := config.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	result := config.DB.Save(&movie)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// DeleteMovie 删除电影
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.Movie{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}

// GetAdminMovies 获取电影列表（管理后台）
func GetAdminMovies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var movies []models.Movie
	var total int64

	offset := (page - 1) * pageSize

	// 获取总记录数
	config.DB.Model(&models.Movie{}).Count(&total)

	// 获取分页数据
	result := config.DB.Offset(offset).Limit(pageSize).Find(&movies)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"results":     movies,
	})
}
