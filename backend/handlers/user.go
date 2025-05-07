package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Estella0129/theater/backend/config"
	"github.com/Estella0129/theater/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthMiddleware JWT验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user.Password = string(hashedPassword)

	// 创建用户
	result := config.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.username") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		} else if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + result.Error.Error()})
		return
	}

	// 清除密码后返回用户信息
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// LoginUser 用户登录
func LoginUser(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}
	//查询用户
	var user models.User
	result := config.DB.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 检查用户是否被冻结
	if user.IsFrozen {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户已被冻结，请联系管理员"})
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成JWT token
	Role := 0
	userID := 0
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 清除密码后返回用户信息
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取总记录数
	config.DB.Model(&models.User{}).Count(&total)

	// 获取分页数据
	result := config.DB.Select("id, username, name, email, role, gender, created_at, updated_at").Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"results":     users,
	})
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := config.DB.Select("id, username, name, email, role, gender, created_at, updated_at").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ToggleFreezeUser 切换用户冻结状态
func ToggleFreezeUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 切换冻结状态
	user.IsFrozen = !user.IsFrozen
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User status updated successfully",
		"is_frozen": user.IsFrozen,
	})
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Gender   string `json:"gender"`
		IsFrozen bool   `json:"is_frozen"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	// 更新用户信息（不改变冻结状态）
	user.Name = updateData.Name
	user.Email = updateData.Email
	user.Gender = updateData.Gender

	result := config.DB.Model(&user).Updates(models.User{
		Username: updateData.Username,
		Name:     updateData.Name,
		Email:    updateData.Email,
		Role:     updateData.Role,
		Gender:   updateData.Gender,
	})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.username") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		} else if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdatePassword(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var passwordData struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password data"})
		return
	}

	// 验证两次新密码是否一致
	if passwordData.NewPassword != passwordData.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New passwords do not match"})
		return
	}

	// 验证当前密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordData.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// CreateUser 管理员创建用户
func CreateUser(c *gin.Context) {
	// 验证当前用户是否为管理员
	// TODO: 从JWT token中获取当前用户角色
	// currentUserRole := getUserRoleFromToken(c)
	// if currentUserRole != "admin" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create users"})
	// 	return
	// }

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user.Password = string(hashedPassword)

	// 创建用户
	result := config.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.username") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		} else if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + result.Error.Error()})
		return
	}

	// 清除密码后返回用户信息
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}
