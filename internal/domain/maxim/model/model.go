package model

import (
	"gorm.io/gorm"
)

type Maxim struct {
	gorm.Model
	Content string `gorm:"not null;size:60"`
	Author  string `gorm:"size:30 not null"`
	Color   string `gorm:"size:10;default:'#ffffff'"`
}

type CreateReq struct {
	Content string `json:"content" binding:"required,max=60"`
	Author  string `json:"author" binding:"required,max=60"`
	Color   string `json:"color"`
}

type UpdateReq struct {
	Content string `json:"content" binding:"required,max=60"`
	Author  string `json:"author" binding:"required,max=60"`
	Color   string `json:"color"`
}

type ListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,max=15"`
	Keyword  string `form:"keyword"`
}

type ListRes struct {
	List  []MaximDTO `json:"list"`
	Pages int        `json:"pages"`
}

type AllRes struct {
	List []MaximDTO `json:"list"`
}
