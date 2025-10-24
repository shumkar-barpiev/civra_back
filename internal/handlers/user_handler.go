package handlers

import (
	"net/http"
	"strconv"

	"civra_back/internal/database"
	"civra_back/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser handles user registration
func CreateUser(c *gin.Context) {
	var userReq models.UserCreateRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", userReq.Email, userReq.Username).First(&existingUser).Error; err == nil {
		Error(c, http.StatusConflict, "User with this email or username already exists")
		return
	}

	user := models.User{
		Username: userReq.Username,
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		Error(c, http.StatusInternalServerError, "Could not create user")
		return
	}

	Created(c, user.ToResponse(), "User created successfully")
}


// LoginUser handles user authentication
func LoginUser(c *gin.Context) {
	var loginReq models.UserLoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if !user.CheckPassword(loginReq.Password) {
		Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	Success(c, user.ToResponse(), "Login successful")
}

// GetUser retrieves a user by ID
func GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, http.StatusNotFound, "User not found")
			return
		}
		Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	Success(c, user.ToResponse(), "User retrieved successfully")
}

// GetUsers retrieves all users (for testing)
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		Error(c, http.StatusInternalServerError, "Could not fetch users")
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	Success(c, userResponses, "Users fetched successfully")
}
