package model

import (
	"blog/internal/response"
	"gorm.io/gorm"
	"time"
)

type Label struct {
	gorm.Model
	Name         string  `gorm:"uniqueIndex;not null;size:30"`
	Introduction *string `gorm:"size:60"`
}

type CreateReq struct {
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"max=60"`
}

type CreateRes struct {
	Code response.Code `json:"-"`
}

type UpdateReq struct {
	ID           uint    `uri:"id" binding:"required"`
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"max=60"`
}

type UpdateRes struct {
	Code response.Code `json:"-"`
}

type DeleteReq struct {
	ID uint `uri:"id" binding:"required"`
}

type DeleteRes struct {
	Code response.Code `json:"-"`
}
type ListReq struct {
	Page uint          `form:"page" binding:"required"`
	Size uint          `form:"size" binding:"required,max=15"`
	Code response.Code `json:"-"`
}

type LabelItem struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Introduction string    `json:"introduction"`
	CreatedTime  time.Time `json:"created_time"`
}
type ListRes struct {
	List  []LabelItem `json:"list"`
	Pages uint        `json:"pages"`
}
