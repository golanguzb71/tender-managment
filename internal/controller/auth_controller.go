package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tender-managment/internal/service"
	"tender-managment/internal/utils"
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
// @Param payload body struct{ Username string `json:"username"`, Email string `json:"email"`, Password string `json:"password"`, Role string `json:"role"`} true "User Registration Payload"
// @Success 201 {object} utils.Response{message=string, data=object{token=string}}
// @Failure 400 {object} utils.Response{message=string, data=object{}}
// @Router /register [post]
func Register(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&payload); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	token, err := authService.RegisterUser(payload.Username, payload.Email, payload.Password, payload.Role)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.Response(c, http.StatusCreated, "User registered successfully", gin.H{"token": token})
}

// Login godoc
// @Summary Login an existing user
// @Description Authenticates a user and returns an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body struct{ Username string `json:"username"`, Password string `json:"password"`} true "User Login Payload"
// @Success 200 {object} utils.Response{message=string, data=object{token=string}}
// @Failure 400 {object} utils.Response{message=string, data=object{}}
// @Failure 401 {object} utils.Response{message=string, data=object{}}
// @Router /login [post]
func Login(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&payload); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	token, err := authService.AuthenticateUser(payload.Username, payload.Password)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.Response(c, http.StatusOK, "Login successful", gin.H{"token": token})
}
