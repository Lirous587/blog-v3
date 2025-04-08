package model

type LabelDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Introduction *string `json:"introduction"`
	EssayCount   uint    `json:"essay_count"`
}

func (l *Label) ConvertToDTO() *LabelDTO {
	return &LabelDTO{
		ID:           l.ID,
		Name:         l.Name,
		Introduction: l.Introduction,
	}
}
