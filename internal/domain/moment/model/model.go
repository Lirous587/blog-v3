package model

import (
	"gorm.io/gorm"
)

type Moment struct {
	gorm.Model
	Title    string `gorm:"size:60;not null"`
	Content  string `gorm:"not null"`
	Location string `gorm:"size:100"`
}

type CreateReq struct {
	Title    string `json:"title" binding:"required,max=60"`
	Content  string `json:"content" binding:"required"`
	Location string `json:"location" binding:"max=100"`
}

type UpdateReq struct {
	Title    string `json:"title" binding:"required,max=60"`
	Content  string `json:"content" binding:"required"`
	Location string `json:"location" binding:"max=100"`
}

type ListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,max=15"`
	Keyword  string `form:"keyword"`
}

type ListRes struct {
	List  []MomentDTO `json:"list"`
	Pages int         `json:"pages"`
}
