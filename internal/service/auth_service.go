package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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
	user, err := as.repo.GetUserByEmail(email)
	if err != nil && err.Error() != "user not found" {
		return "", err
	}

	if user != nil {
		return "", errors.New("email already exists")
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

func (as *AuthService) AuthenticateUser(username, password string) (string, error) {
	user, err := as.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
