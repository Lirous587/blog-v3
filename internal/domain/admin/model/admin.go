package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null;size:255"`
	HashPassword []byte `gorm:"type:varchar(255);not null"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=70"`
}

// LoginRequest 登录请求
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	ID           uint   `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtPaylod struct {
	ID uint `json:"id"`
}
