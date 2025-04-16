package handlers

import (
	"net/http"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
)

// GetGenres 获取所有电影类型
func GetGenres(c *gin.Context) {
	var genres []models.Genre
	if err := config.DB.Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取电影类型失败"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// CreateGenre 创建新的电影类型
func CreateGenre(c *gin.Context) {
	var newGenre models.Genre
	if err := c.ShouldBindJSON(&newGenre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数解析失败"})
		return
	}

	if err := config.DB.Create(&newGenre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建电影类型失败"})
		return
	}

	c.JSON(http.StatusCreated, newGenre)
}

// UpdateGenre 更新电影类型
func UpdateGenre(c *gin.Context) {
	id := c.Param("id")
	var genre models.Genre
	if err := config.DB.First(&genre, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "电影类型未找到"})
		return
	}

	var updatedGenre models.Genre
	if err := c.ShouldBindJSON(&updatedGenre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数解析失败"})
		return
	}

	if err := config.DB.Model(&genre).Updates(updatedGenre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新电影类型失败"})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// DeleteGenre 删除电影类型
func DeleteGenre(c *gin.Context) {
	id := c.Param("id")
	var genre models.Genre
	if err := config.DB.First(&genre, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "电影类型未找到"})
		return
	}

	if err := config.DB.Delete(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除电影类型失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "电影类型删除成功"})
}

// GetAdminGenres 获取管理员可见的所有电影类型
func GetAdminGenres(c *gin.Context) {
	var genres []models.Genre
	if err := config.DB.Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取管理员可见电影类型失败"})
		return
	}

	c.JSON(http.StatusOK, genres)
}
