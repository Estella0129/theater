package handlers

import (
	"net/http"
	"strconv"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
)

// GetMovies 获取电影列表，支持分页
func GetMovies(c *gin.Context) {
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

// GetMovie 获取单个电影详情
func GetMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	result := config.DB.First(&movie, id)
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

// SearchMovies 搜索电影
func SearchMovies(c *gin.Context) {
	query := c.Query("query")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var movies []models.Movie
	var total int64

	offset := (page - 1) * pageSize

	// 获取总记录数
	db := config.DB.Model(&models.Movie{})
	if query != "" {
		db = db.Where("title LIKE ?", "%"+query+"%")
	}
	db.Count(&total)

	// 获取分页数据
	result := db.Offset(offset).Limit(pageSize).Find(&movies)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search movies"})
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
