package model

type RegisterModel struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
