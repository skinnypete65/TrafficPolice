package dto

type Camera struct {
	ID           string  `json:"camera_id,omitempty"`
	CameraTypeID string  `json:"camera_type_id,omitempty" validate:"required,uuid"`
	Latitude     float64 `json:"latitude,omitempty" validate:"required"`
	Longitude    float64 `json:"longitude,omitempty" validate:"required"`
	ShortDesc    string  `json:"short_desc,omitempty" validate:"required"`
}

type CameraType struct {
	ID   string `json:"camera_id,omitempty"`
	Name string `json:"camera_name,omitempty" validate:"required"`
}

type RegisterCamera struct {
	Camera Camera `json:"camera" validate:"required"`
	SignUp SignUp `json:"sign_up" validate:"required"`
}
