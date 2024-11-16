package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
)

var (
	authService *service.AuthService
)

func SetAuthService(authSer *service.AuthService) {
	authService = authSer
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user and returns an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.RegisterModel true "User Registration Payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	var payload model.RegisterModel
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	token, err := authService.RegisterUser(payload.Username, payload.Email, payload.Password, payload.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

// Login godoc
// @Summary Login an existing user
// @Description Authenticates a user and returns an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.LoginModel true "User Login Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	var payload model.LoginModel

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	token, status, err := authService.AuthenticateUser(payload.Username, payload.Password)
	if err != nil {
		c.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
