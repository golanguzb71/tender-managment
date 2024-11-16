package controller

import "tender-managment/internal/service"

var (
	userService *service.UserService
)

func SetUserService(userSer *service.UserService) {
	userService = userSer
}
