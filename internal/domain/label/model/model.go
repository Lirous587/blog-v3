package model

import (
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Name         string  `gorm:"not null;size:30"`
	Introduction *string `gorm:"size:60"`
}

type CreateReq struct {
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"omitempty,max=60"`
}

type UpdateReq struct {
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"omitempty,max=60"`
}

type ListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,max=15"`
	Keyword  string `form:"keyword"`
}

type ListRes struct {
	List  []LabelDTO `json:"list"`
	Pages int        `json:"pages"`
}

type AllRes struct {
	List []LabelDTO `json:"list"`
}
