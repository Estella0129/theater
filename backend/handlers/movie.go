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

// SearchMovies 通过关键字搜索电影
func SearchMovies(c *gin.Context) {
	// 参数解析
	keyword := c.Query("query")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20

	// 验证参数
	if len(keyword) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "搜索关键字至少需要2个字符",
		})
		return
	}

	// 构造查询
	query := config.DB.Model(&models.Movie{}).
		Where("title LIKE ? OR description LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var movies []models.Movie
	result := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&movies)

	// 错误处理
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库查询失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"data": movies,
		"meta": gin.H{
			"current_page": page,
			"per_page":     pageSize,
			"total":        total,
		},
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
