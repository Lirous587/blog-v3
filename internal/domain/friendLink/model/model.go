package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPublished = "published"
	StatusPending   = "pending"
)

type FriendLink struct {
	gorm.Model
	Introduction string `gorm:"not null;size:80"`
	SiteName     string `gorm:"not null;size:80"`
	Url          string `gorm:"not null;size:120;uniqueIndex"`
	Logo         string `gorm:"not null;size:120"`
	Status       Status `gorm:"default:published;check status IN ('published','pending')"`
	Email        string `gorm:"not null;size:80"`
}

func (fl *FriendLink) BeforeSave(tx *gorm.DB) error {
	if fl.Status != StatusPublished && fl.Status != StatusPending {
		if fl.Status == "" {
			fl.Status = StatusPending
			return nil
		}
		return errors.New("status 必须为 'published' 或 'pending'")
	}
	return nil
}

type CreateReq struct {
	Introduction string `json:"introduction" binding:"required,max=80"`
	SiteName     string `json:"site_name" binding:"required,max=80"`
	Url          string `json:"url" binding:"required,domain_url,max=120"`
	Logo         string `json:"logo" binding:"required,url,max=120"`
	Email        string `json:"email" binding:"required,email,max=80"`
}

type UpdateReq struct {
	Introduction string `json:"introduction" binding:"required,max=80"`
	SiteName     string `json:"site_name" binding:"required,max=80"`
	Url          string `json:"url" binding:"required,domain_url,max=120"`
	Logo         string `json:"logo" binding:"required,url,max=120"`
	Status       Status `json:"status" binding:"oneof=published pending"`
	Email        string `json:"email" binding:"required,email,max=80"`
}

type UpdateStatusReq struct {
	Status Status `json:"status" binding:"oneof=published pending"`
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

type ApplyReq struct {
}
