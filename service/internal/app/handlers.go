package app

import (
	"TrafficPolice/internal/camera"
	"TrafficPolice/internal/transport/rabbitmq"
	"TrafficPolice/internal/transport/rest"
	"github.com/go-playground/validator/v10"
)

type handlers struct {
	rating      *rest.RatingHandler
	auth        *rest.AuthHandler
	camera      *rest.CameraHandler
	caseHandler *rest.CaseHandler
	contactInfo *rest.ContactInfoHandler
	violation   *rest.ViolationHandler
	expert      *rest.ExpertHandler
	training    *rest.TrainingHandler
	director    *rest.DirectorHandler
}

func newHandlers(
	s *services,
	c *converters,
	validate *validator.Validate,
	finePublisher *rabbitmq.FinePublisher,
) *handlers {
	cameraParser := camera.NewParser(s.camera)
	return &handlers{
		rating:      rest.NewRatingHandler(s.rating, c.rating),
		auth:        rest.NewAuthHandler(s.auth, validate, c.userInfo, c.auth),
		camera:      rest.NewCameraHandler(s.camera, s.auth, validate, c.camera),
		caseHandler: rest.NewCaseHandler(s.caseService, s.img, s.camera, c.caseConverter, cameraParser),
		contactInfo: rest.NewContactInfoHandler(s.contactInfo),
		violation:   rest.NewViolationHandler(s.violation),
		expert: rest.NewExpertHandler(
			s.img, s.expert, s.rating, finePublisher, c.caseConverter, c.caseDecision,
		),
		training: rest.NewTrainingHandler(
			s.training, s.pagination, validate, c.caseConverter, c.pagination, c.solvedCases,
		),
		director: rest.NewDirectorHandler(s.director, c.caseConverter, c.analytics),
	}
}
