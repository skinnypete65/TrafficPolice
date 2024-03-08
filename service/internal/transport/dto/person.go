package dto

type Person struct {
	ID       string `json:"ID,omitempty"`
	PhoneNum string `json:"phone_num,omitempty"`
	Email    string `json:"email,omitempty"`
	VkID     string `json:"vk_id,omitempty"`
	TgID     string `json:"tg_id,omitempty"`
}
