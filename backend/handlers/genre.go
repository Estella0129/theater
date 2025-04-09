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
