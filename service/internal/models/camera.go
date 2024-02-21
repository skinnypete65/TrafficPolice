package models

type Camera struct {
	ID           string  `json:"camera_id"`
	CameraTypeID string  `json:"camera_type_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	ShortDesc    string  `json:"short_desc"`
}

type CameraType struct {
	ID   string `json:"camera_id"`
	Name string `json:"camera_name" validate:"required"`
	Desc string `json:"camera_desc" validate:"required"`
}
