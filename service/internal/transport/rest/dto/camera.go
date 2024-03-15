package dto

type Camera struct {
	ID           string  `json:"camera_id,omitempty"`
	CameraTypeID string  `json:"camera_type_id,omitempty" validate:"required,uuid"`
	Latitude     float64 `json:"latitude,omitempty" validate:"required"`
	Longitude    float64 `json:"longitude,omitempty" validate:"required"`
	ShortDesc    string  `json:"short_desc,omitempty" validate:"required"`
}

type CameraTypeIn struct {
	Name string `json:"camera_name,omitempty" validate:"required"`
}

type CameraIn struct {
	CameraTypeID string  `json:"camera_type_id,omitempty" validate:"required,uuid"`
	Latitude     float64 `json:"latitude,omitempty" validate:"required"`
	Longitude    float64 `json:"longitude,omitempty" validate:"required"`
	ShortDesc    string  `json:"short_desc,omitempty" validate:"required"`
}

type RegisterCamera struct {
	CameraIn CameraIn `json:"camera" validate:"required"`
	SignUp   SignUp   `json:"sign_up" validate:"required"`
}
