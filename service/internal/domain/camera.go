package domain

type Camera struct {
	ID         string
	CameraType CameraType
	Latitude   float64
	Longitude  float64
	ShortDesc  string
}

type CameraType struct {
	ID   string
	Name string
}
