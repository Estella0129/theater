package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
)

// GetPeople 获取人物列表，支持分页和搜索
func GetPeoples(c *gin.Context) {
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

	var people []models.People
	var total int64

	offset := (page - 1) * pageSize

	dbQuery := config.DB.Model(&models.People{})
	if searchQuery != "" {
		dbQuery = dbQuery.Where("name LIKE ? OR original_name LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// 获取总记录数
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取人物总数失败"})
		return
	}

	// 获取分页数据
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&people).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取人物列表失败"})
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
		"results":     people,
	})
}

// GetPeople 获取单个人物详情
func GetPeople(c *gin.Context) {
	id := c.Param("id")

	var People models.People
	result := config.DB.Preload("Credits").Preload("Credits.Movie").First(&People, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "People not found"})
		return
	}

	c.JSON(http.StatusOK, People)
}

// CreatePeople 创建人物
func CreatePeople(c *gin.Context) {
	var People models.People
	if err := c.ShouldBindJSON(&People); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result := config.DB.Create(&People)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create People"})
		return
	}

	c.JSON(http.StatusCreated, People)
}

// UpdatePeople 更新人物信息
func UpdatePeople(c *gin.Context) {
	id := c.Param("id")

	var People models.People
	if err := config.DB.First(&People, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "People not found"})
		return
	}

	if err := c.ShouldBindJSON(&People); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	result := config.DB.Save(&People)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update People"})
		return
	}

	c.JSON(http.StatusOK, People)
}

// DeletePeople 删除人物
func DeletePeople(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.People{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete People"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "People deleted successfully"})
}

// GetAdminPeople 获取人物列表（管理后台）
func GetAdminPeople(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	searchQuery := c.Query("search")

	var people []models.People
	var total int64

	offset := (page - 1) * pageSize

	db := config.DB
	if searchQuery != "" {
		db = db.Where("name LIKE ?", "%"+searchQuery+"%")
	}

	// 获取总记录数
	db.Model(&models.People{}).Count(&total)

	// 获取分页数据
	result := db.Offset(offset).Limit(pageSize).Find(&people)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch people"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"results":     people,
	})
}
