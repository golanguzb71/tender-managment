package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (ur *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password, role FROM users WHERE email = $1`
	var user User
	row := ur.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password, role FROM users WHERE username = $1`
	var user User
	row := ur.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error fetching user: %v", err)
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}

func (ur *UserRepository) CreateUser(username, email, password, role string) (int, error) {
	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	var userID int
	err := ur.db.QueryRow(query, username, email, password, role).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}
	return userID, nil
}
