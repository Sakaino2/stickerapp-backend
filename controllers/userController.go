package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"jorge.com/inventari_app_backend/initializers"
	"jorge.com/inventari_app_backend/models"
	"jorge.com/inventari_app_backend/utils"
)

func UpdateUser(context *gin.Context) {
	var body struct {
		Username string
		Email    string
		Password string
		RolID    int
	}

	context.Bind(&body)

	id := context.Param("id")

	var user models.User
	initializers.DATABASE.First(&user, id)

	initializers.DATABASE.Model(&user).Updates(models.User{
		Username: body.Username,
	})

	context.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func DeleteUser(context *gin.Context) {
	id := context.Param("id")

	initializers.DATABASE.Delete(&models.User{}, id)

	context.Status(http.StatusOK)
}

func Signup(context *gin.Context) {
	var body struct {
		Username string
		Email    string
	}

	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Email), 10)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Username: strings.ToUpper(body.Username),
		Email:    strings.ToUpper(body.Email),
		Password: string(hash),
	}

	result := initializers.DATABASE.Create(&user)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(
		http.StatusCreated,
		gin.H{
			"message": "User created successfully.",
		},
	)
}

func Login(context *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User

	initializers.DATABASE.First(&user, "email = ?", strings.ToUpper(body.Email))

	if user.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Validate credentials
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token, refreshToken, err := utils.GenerateToken(strings.ToUpper(body.Email), user.ID)

	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//Response
	context.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh-token": refreshToken,
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
	})
}

func RefreshUserToken(context *gin.Context) {
	var body struct {
		refreshToken string
	}

	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	newToken, newRefreshToken, err := utils.RefreshToken(body.refreshToken)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token":         newToken,
		"refresh-token": newRefreshToken,
	})
}
