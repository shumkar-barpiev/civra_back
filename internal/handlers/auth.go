package handlers

import (
	"net/http"
	"os"
	"time"

	"civra_back/internal/database"
	"civra_back/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // set in .env

// LoginUser handles user authentication and token generation
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

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 1-day expiry
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user data
	Success(c, gin.H{
		"token": tokenString,
		"user":  user.ToResponse(),
	}, "Login successful")
}
