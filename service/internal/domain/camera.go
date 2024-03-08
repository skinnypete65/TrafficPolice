package domain

type Camera struct {
	ID           string  `json:"camera_id"`
	CameraTypeID string  `json:"camera_type_id" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	ShortDesc    string  `json:"short_desc" validate:"required"`
}

type CameraType struct {
	ID   string `json:"camera_id"`
	Name string `json:"camera_name" validate:"required"`
}
