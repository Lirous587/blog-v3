package model

import (
	"blog/internal/domain/label/model"
	"blog/utils"
)

type EssayDTO struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Introduction *string          `json:"introduction,omitempty"`
	CreatedAt    string           `json:"create_at"`
	Content      string           `json:"content,omitempty"`
	PreviewTheme string           `json:"preview_theme"`
	CodeTheme    string           `json:"code_theme"`
	ImgUrl       *string          `json:"img_url"`
	Labels       []model.LabelDTO `json:"labels"`
	Priority     int8             `json:"priority"`
}

func (essay *Essay) ConvertToDTO() *EssayDTO {
	if essay == nil {
		return nil
	}
	var labels []model.LabelDTO
	if essay.Labels != nil {
		labels = make([]model.LabelDTO, len(essay.Labels))
		for i, label := range essay.Labels {
			labels[i] = *label.ConvertToDTO()
		}
	}

	return &EssayDTO{
		ID:           essay.ID,
		Name:         essay.Name,
		Introduction: essay.Introduction,
		Content:      essay.Content,
		Priority:     essay.Priority,
		PreviewTheme: essay.PreviewTheme,
		CodeTheme:    essay.CodeTheme,
		ImgUrl:       essay.ImgUrl,
		CreatedAt:    utils.FormatTime(essay.CreatedAt),
		Labels:       labels,
	}
}
