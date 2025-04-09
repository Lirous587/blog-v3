package model

type MaximDTO struct {
	ID           uint   `json:"id"`
	Introduction string `json:"introduction"`
	SiteName     string `json:"siteName"`
	Url          string `json:"url"`
	Logo         string `json:"logo"`
	Status       Status `json:"status"`
}

func (fl *FriendLink) ConvertToDTO() *MaximDTO {
	return &MaximDTO{
		ID:           fl.ID,
		Introduction: fl.Introduction,
		SiteName:     fl.SiteName,
		Url:          fl.Url,
		Logo:         fl.Logo,
		Status:       fl.Status,
	}
}
