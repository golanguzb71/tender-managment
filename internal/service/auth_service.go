package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	repository "tender-managment/internal/db/repo"
	"tender-managment/internal/utils"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (as *AuthService) RegisterUser(username, email, password, role string) (string, error) {
	if username == "" || email == "" {
		return "", errors.New("username or email cannot be empty")
	}

	if strings.Contains(email, "-") && !strings.Contains(email, "@gmail.com") {
		return "", errors.New("invalid email format")
	}

	user, _ := as.repo.GetUserByEmail(email)

	if user != nil {
		return "", errors.New("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userId, err := as.repo.CreateUser(username, email, string(hashedPassword), role)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(userId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) AuthenticateUser(username, password string) (string, int, error) {
	if username == "" || password == "" {
		return "", http.StatusBadRequest, errors.New("Username and password are required")
	}

	user, err := as.repo.GetUserByUsername(username)
	if err != nil {
		return "", http.StatusNotFound, errors.New("User not found")
	}

	err = utils.ComparePasswords(user.Password, password)
	if err != nil {
		return "", http.StatusUnauthorized, errors.New("Invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return token, http.StatusOK, nil
}
