package model

type MaximDTO struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Color   string `json:"color"`
}

func (m *Maxim) ConvertToDTO() *MaximDTO {
	return &MaximDTO{
		ID:      m.ID,
		Content: m.Content,
		Author:  m.Author,
		Color:   m.Color,
	}
}
