package model

import (
	"blog/internal/domain/label/model"
	"gorm.io/gorm"
)

type Essay struct {
	gorm.Model
	Name         string        `gorm:"not null;size:30;index:idx_essay_name"`
	Introduction *string       `gorm:"size:160"`
	Content      string        `gorm:"type:text;not null"`
	PreviewTheme string        `gorm:"not null;size:20"`
	CodeTheme    string        `gorm:"not null;size:20"`
	ImgUrl       *string       `gorm:"size:512;default:null;comment:文章封面图片链接"`
	Labels       []model.Label `gorm:"many2many:essay_labels;"` // 多对多关系
	VisitedTimes uint          `gorm:"default:1"`
	Priority     uint8         `gorm:"default:50;index:idx_essay_priority"`
}

type CreateReq struct {
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"omitempty,max=60"`
	Content      string  `json:"content" binding:"required"`
	Priority     uint8   `json:"priority" binding:"required,min=0,max=100"`
	PreviewTheme string  `json:"preview_theme" binding:"required,max=20"`
	CodeTheme    string  `json:"code_theme" binding:"required,max=20"`
	ImgUrl       *string `json:"img_url" binding:"url,max=512"`
	LabelIDs     []uint  `json:"label_ids" binding:"unique,min=1,dive,gt=0"`
}

type UpdateReq struct {
	Name         string  `json:"name" binding:"required,max=30"`
	Introduction *string `json:"introduction" binding:"omitempty,max=60"`
	Content      string  `json:"content" binding:"required"`
	Priority     uint8   `json:"priority" binding:"required,min=0,max=100"`
	PreviewTheme string  `json:"preview_theme" binding:"required,max=20"`
	CodeTheme    string  `json:"code_theme" binding:"required,max=20"`
	ImgUrl       *string `json:"img_url" binding:"url,max=512"`
	LabelIDs     []uint  `json:"label_ids" binding:"required,unique,min=1,dive,gt=0"`
}

type ListReq struct {
	Page     int  `form:"page" binding:"required,min=1"`
	PageSize int  `form:"page_size" binding:"required,max=15"`
	LabelID  uint `form:"label_id" binding:"gt=0"`
}

type ListRes struct {
	List  []EssayDTO `json:"list"`
	Pages int        `json:"pages"`
}

type Timeline struct {
	Data    string     `json:"data"`
	Records []EssayDTO `json:"records"`
}
