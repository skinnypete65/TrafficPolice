package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type CameraConverter struct {
}

func NewCameraConverter() *CameraConverter {
	return &CameraConverter{}
}

func (c *CameraConverter) MapRegisterCameraDtoToDomain(camera dto.CameraIn,
	signUp dto.SignUp,
) domain.RegisterCamera {
	return domain.RegisterCamera{
		Camera: domain.Camera{
			ID:         "",
			CameraType: domain.CameraType{ID: camera.CameraTypeID},
			Latitude:   camera.Latitude,
			Longitude:  camera.Longitude,
			ShortDesc:  camera.ShortDesc,
		},
		Username: signUp.Username,
		Password: signUp.Password,
	}
}
