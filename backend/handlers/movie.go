package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"

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

	dbQuery.Preload("Director", "job = ?", "Director").Preload("Director.People")
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
	result := config.DB.
		Preload("Credits").
		Preload("Credits.People").
		Preload("Genres").
		Preload("Images").First(&movie, id)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Create(&movie).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建电影失败"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie 更新电影信息
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的更新数据: %s", err.Error())})
		return
	}

	tx := config.DB.Begin()
	//if err := tx.Preload("Genres").First(&movie, id).Error; err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "电影未找到"})
	//	return
	//}

	// 删除原有关联数据
	if err := tx.Where("movie_id = ?", movie.ID).Delete(&models.Credit{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清除演职人员失败"})
		return
	}

	// 删除原有关联数据
	if err := tx.Where("movie_id = ?", movie.ID).Delete(&models.MovieImage{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清除图片失败"})
		return
	}

	imagePathList := []string{}
	for _, image := range movie.Images {
		imagePathList = append(imagePathList, image.FilePath)
	}

	// 删除原有关联数据
	if err := tx.Where("file_path in ?", imagePathList).Delete(&models.Image{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清除图片失败"})
		return
	}

	if err := tx.Model(&movie).Where(id).Updates(movie).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新电影失败"})
		return
	}
	//
	//// 更新类型关联
	//var genres []models.Genre
	//if err := tx.Where("id IN ?", movie.Genres).Find(&genres).Error; err != nil {
	//	tx.Rollback()
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的电影类型"})
	//	return
	//}
	//tx.Model(&movie).Association("Genres").Replace(genres)
	//
	//
	//// 重新创建演职人员
	//for _, credit := range movie.Credits {
	//	credit.MovieID = int(movie.ID)
	//	if err := tx.Create(&credit).Error; err != nil {
	//		tx.Rollback()
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存演职人员失败"})
	//		return
	//	}
	//}

	tx.Commit()
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

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "文件上传失败",
			"message": err.Error(),
		})
		return
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "不支持的图片格式",
			"detail": file.Header.Get("Content-Type"),
		})
		return
	}

	// 生成唯一文件名
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	dstPath := filepath.Join("images", fileName)

	// 设置CORS头
	c.Header("Access-Control-Allow-Origin", "*")
	// 保存文件
	if error := c.SaveUploadedFile(file, dstPath); error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "文件保存失败",
			"message": error.Error(),
		})
		return
	}

	// 获取图片元数据
	dstFile, err := os.Open(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "打开保存文件失败",
			"message": err.Error(),
		})
		return
	}
	defer dstFile.Close()

	// 重置文件读取位置
	dstFile.Seek(0, 0)
	img, err := imaging.Open(dstFile.Name(), imaging.AutoOrientation(true))
	if err != nil {
		log.Printf("解码失败 文件路径:%s 错误详情:%v", dstPath, err)
		_ = os.Remove(dstPath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "解析图片元数据失败",
			"message": fmt.Sprintf("文件格式可能损坏，详情：%s", err.Error()),
		})
		return
	}

	// 移除未使用的str变量
	fmt.Println("图片元数据获取错误:", err)

	c.JSON(http.StatusOK, gin.H{
		"file_path":    "/" + fileName,
		"width":        img.Bounds().Dx(),
		"height":       img.Bounds().Dy(),
		"aspect_ratio": float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()),
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
