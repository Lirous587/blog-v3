package model

type LabelDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Introduction *string `json:"introduction,omitempty"`
}

func (l *Label) ConvertToDTO() *LabelDTO {
	return &LabelDTO{
		ID:           l.ID,
		Name:         l.Name,
		Introduction: l.Introduction,
	}

	//dtos := make([]model.LabelDTO, len(labels))
	//for i, label := range labels {
	//	dtos[i] = model.LabelDTO{
	//		ID:           label.ID,
	//		Name:         label.Name,
	//		Introduction: label.Introduction,
	//	}
	//}
}
