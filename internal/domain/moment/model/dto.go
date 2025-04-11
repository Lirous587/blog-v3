package model

import "time"

type MomentDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Location  string    `json:"location"`
}

func (m *Moment) ConvertToDTO() *MomentDTO {
	return &MomentDTO{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		Content:   m.Content,
		Title:     m.Title,
		Location:  m.Location,
	}
}
