package models

type Camera struct {
	ID         string
	CameraType CameraType
	latitude   float64
	longitude  float64
	ShortDesc  string
}

type CameraType struct {
	ID   string `json:"camera_id"`
	Name string `json:"camera_name" validate:"required"`
	Desc string `json:"camera_desc" validate:"required"`
}
