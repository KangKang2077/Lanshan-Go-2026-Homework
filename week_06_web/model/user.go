package model

import (
	"time"
)

type User struct {
	Id           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Create_at    time.Time `json:"create_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=1,max=12"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=1,max=12"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}
